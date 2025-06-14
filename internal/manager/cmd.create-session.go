package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/letstrygo/letstry/internal/config"
	"github.com/letstrygo/letstry/internal/config/editors"
	"github.com/letstrygo/letstry/internal/environment"
	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/util/identifier"
	"github.com/otiai10/copy"
)

type CreateSessionArguments struct {
	Source string `json:"source"`
}

func (s *manager) CreateSession(ctx context.Context, args CreateSessionArguments) (Session, error) {
	var (
		zeroValue Session
		err       error
	)

	sourceType, err := s.GetSessionSourceType(ctx, args.Source)
	if err != nil {
		return zeroValue, err
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return zeroValue, err
	}

	editor, err := cfg.GetDefaultEditor()
	if err != nil {
		return zeroValue, err
	}

	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return zeroValue, err
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "letstry")
	if err != nil {
		return zeroValue, fmt.Errorf("failed to create temporary directory: %v", err)
	}

	logger.Printf("found source type: %s\n", sourceType.FormattedString())

	// Handle "With" arguments
	switch sourceType {
	case SessionSourceTypeBlank:
		// Do nothing
	case SessionSourceTypeDirectory:
		err = s.fillWorkspaceFromDirectory(ctx, args.Source, tempDir)
	case SessionSourceTypeRepository:
		err = s.fillWorkspaceFromRepository(ctx, args.Source, tempDir)
	case SessionSourceTypeTemplate:
		err = s.fillWorkspaceFromTemplate(ctx, args.Source, tempDir)
	}

	if err != nil {
		return zeroValue, err
	}

	// Launch the editor
	cmd, err := s.launchEditor(ctx, editor, tempDir)
	if err != nil {
		return zeroValue, err
	}

	// Create the session
	session, err := s.persistSession(ctx, sourceType, cmd, editor, args.Source, tempDir)
	if err != nil {
		return zeroValue, err
	}

	// Monitor session
	err = s.monitor(ctx, session)

	return *session, err
}

func (s *manager) monitor(ctx context.Context, session *Session) error {
	appEnvironment, err := environment.EnvironmentFromContext(ctx)
	if err != nil {
		return err
	}

	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	if appEnvironment.DebuggerAttached {
		logger.Printf("skipping monitor process for session %s (debugger attached)\n", session.FormattedID())
		logger.Printf("starting monitor session for session %s with delay %v\n", session.FormattedID(), session.Editor.ProcessCaptureDelay)
		err = s.MonitorSession(ctx, MonitorSessionArguments{
			Delay:        session.Editor.ProcessCaptureDelay,
			TrackingType: session.Editor.TrackingType,
			Location:     session.Location,
			PID:          session.PID,
		})
		if err != nil {
			return err
		}
	} else {
		// Call this application again, but start it in the background as it's own process.
		// This will allow the user to continue using the current terminal session.
		logger.Printf("starting monitor session for session %s with delay %v\n", session.FormattedID(), session.Editor.ProcessCaptureDelay)
		cmd := exec.Command(os.Args[0], "monitor", fmt.Sprintf("%v", session.Editor.ProcessCaptureDelay), session.Location, fmt.Sprintf("%v", session.PID), session.Editor.TrackingType.String())
		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to start monitor process: %v", err)
		}
		logger.Printf("monitor process started with PID %v\n", cmd.Process.Pid)
		logger.Printf("session created: %s\n", session.String())
	}

	return nil
}

func (s *manager) persistSession(ctx context.Context, sourceType SessionSourceType, cmd *exec.Cmd, editor editors.Editor, source string, tempDir string) (*Session, error) {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return nil, err
	}

	session := Session{
		ID:       identifier.NewID(),
		Location: tempDir,
		PID:      cmd.Process.Pid,
		Source:   sessionSource{SourceType: sourceType, Value: source},
		Editor:   editor,
	}

	logger.Printf("persisting session %s\n", session.FormattedID())

	// Save the session
	err = s.addSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *manager) launchEditor(ctx context.Context, editor editors.Editor, tempDir string) (*exec.Cmd, error) {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return nil, err
	}

	logger.Printf("launching editor %s\n", editor.String())
	cfgArgs := strings.Split(editor.Args, " ")
	cmdArgs := append(cfgArgs, tempDir)
	cmd := exec.Command(editor.ExecPath, cmdArgs...)
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to run editor: %v", err)
	}

	return cmd, nil
}

func (s *manager) fillWorkspaceFromTemplate(ctx context.Context, source string, tempDir string) error {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	logger.Printf("reading from template %s\n", source)

	// Check if the specified template exists.
	template, err := s.GetTemplate(ctx, source)
	if err != nil {
		return err
	}

	// Copy the template to the temporary directory
	err = copy.Copy(template.AbsolutePath(ctx), tempDir)
	if err != nil {
		return fmt.Errorf("failed to load template %s: %s", source, err)
	}

	return nil
}

func (s *manager) fillWorkspaceFromRepository(ctx context.Context, source string, tempDir string) error {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	logger.Printf("cloning repository %s\n", source)

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: source,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	return nil
}

func (s *manager) fillWorkspaceFromDirectory(ctx context.Context, source string, tempDir string) error {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	logger.Printf("copying data from source directory %v", source)

	absPath, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		return fmt.Errorf("directory %s does not exist", absPath)
	}

	// Copy the directory to the temporary directory
	err = copy.Copy(absPath, tempDir)
	if err != nil {
		return fmt.Errorf("failed to copy directory: %v", err)
	}

	return nil
}

func (s *manager) addSession(ctx context.Context, sess Session) error {
	sessions, err := s.ListSessions(ctx)
	if err != nil {
		return err
	}

	// check if the session already exists by the same name
	for _, session := range sessions {
		if session.ID == sess.ID {
			return fmt.Errorf("session with ID %s already exists", sess.ID)
		}
	}

	// add the session to the list of sessions
	sessions = append(sessions, sess)

	// save the sessions
	file, err := s.storage.OpenFile("sessions.json")
	if err != nil {
		return fmt.Errorf("failed to open sessions file: %v", err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(sessions, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal sessions: %v", err)
	}

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate sessions: %v", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write sessions: %v", err)
	}

	err = file.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync sessions file: %v", err)
	}

	return nil
}
