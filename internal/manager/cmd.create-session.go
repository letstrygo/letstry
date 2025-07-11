package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/letstrygo/letstry/internal/config"
	"github.com/letstrygo/letstry/internal/config/editors"
	"github.com/letstrygo/letstry/internal/environment"
	"github.com/letstrygo/letstry/internal/util/identifier"
	"github.com/otiai10/copy"
	"github.com/samber/lo"
)

type CreateSessionArguments struct {
	Source             string `json:"source"`
	ForceRequireExport bool   `json:"force_require_export"`
}

func (s *manager) CreateSession(ctx context.Context, args CreateSessionArguments) (*Session, error) {
	var (
		err error
	)

	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	editor, err := cfg.GetDefaultEditor()
	if err != nil {
		return nil, err
	}

	src, err := s.parseSessionSource(ctx, args.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session source: %v", err)
	}

	// Create temporary directory
	projectName := fmt.Sprintf("%v-lt%d", src.ShortValue(), time.Now().Unix())
	storageDir := filepath.Join(cfg.LTPath, projectName)

	requireExport := cfg.LTPath == "" || cfg.RequireExport || args.ForceRequireExport

	if requireExport {
		tempDir, err := os.MkdirTemp("", projectName)
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary directory: %v", err)
		}

		storageDir = tempDir
	} else {
		_, err = os.Stat(storageDir)
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to create project directory: directory already exists")
		}

		err = os.MkdirAll(storageDir, 0660)
		if err != nil {
			return nil, fmt.Errorf("failed to create project directory: %v", err)
		}
	}

	id := identifier.NewID()

	// Fill workspace based on source type.
	err = s.fillWorkspace(ctx, src, storageDir)
	if err != nil {
		return nil, err
	}

	// Launch the editor
	cmd, err := s.launchEditor(ctx, editor, storageDir)
	if err != nil {
		return nil, err
	}

	// Monitor session, automatically purging it from the cache once closed.
	if requireExport {
		// Cache the session in the file system.
		session, err := s.prepareMonitor(ctx, id, cmd, editor, src, storageDir)
		if err != nil {
			return nil, err
		}

		return session, s.monitor(ctx, session)
	}

	return nil, nil
}

func (s *manager) monitor(ctx context.Context, session *Session) error {
	appEnvironment, err := environment.EnvironmentFromContext(ctx)
	if err != nil {
		return err
	}

	if appEnvironment.DebuggerAttached {
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
		cmd := exec.Command(os.Args[0], "monitor", fmt.Sprintf("%v", session.Editor.ProcessCaptureDelay), session.Location, fmt.Sprintf("%v", session.PID), session.Editor.TrackingType.String())
		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to start monitor process: %v", err)
		}
	}

	return nil
}

func (s *manager) parseSessionSource(ctx context.Context, source string) (Source, error) {
	var zeroValue Source

	sourceType, err := s.GetSessionSourceType(ctx, source)
	if err != nil {
		return zeroValue, err
	}

	return Source{sourceType, source}, nil
}

func (s *manager) prepareMonitor(ctx context.Context, id identifier.ID, cmd *exec.Cmd, editor editors.Editor, source Source, storageDir string) (*Session, error) {
	pid, err := s.locatePid(cmd.Process.Pid)
	if err != nil {
		return nil, err
	}

	session := Session{
		ID:       id,
		Location: storageDir,
		PID:      pid,
		Source:   source,
		Editor:   editor,
	}

	// Save the session
	err = s.addSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// locatePid will attempt to locate the PID of the most recent vscode process
// by assuming that vscode processes will have the ps cmd of `/usr/share/code/code`
func (s *manager) locatePid(pid int) (int, error) {
	// This functionality is only supported on linux.
	if runtime.GOOS != "linux" {
		return pid, nil
	}

	time.Sleep(2 * time.Second)

	psCmd := exec.Command(
		"ps",
		"-eo", "pid,etimes,cmd",
		"--sort=etimes",
	)

	grepCmd := exec.Command(
		"grep",
		"/usr/share/code/code",
	)

	// Pipe psCmd's stdout to grepCmd's stdin
	psOut, err := psCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	grepCmd.Stdin = psOut

	// Capture grepCmd's output
	var out bytes.Buffer
	grepCmd.Stdout = &out

	// Start psCmd
	if err := psCmd.Start(); err != nil {
		return pid, err
	}

	// Start grepCmd
	if err := grepCmd.Start(); err != nil {
		return pid, err
	}

	// Wait for psCmd to finish
	if err := psCmd.Wait(); err != nil {
		return pid, err
	}

	// Wait for grepCmd to finish
	if err := grepCmd.Wait(); err != nil {
		return pid, err
	}

	output := out.String()

	// Split output into lines and trim whitespace
	lines := lo.FilterMap(
		strings.Split(output, "\n"),
		func(line string, _ int) (string, bool) {
			trimmed := strings.TrimSpace(line)
			return trimmed, trimmed != ""
		},
	)

	// Get the second line (index 1), if it exists
	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		if len(fields) > 0 {
			_pid, err := strconv.Atoi(fields[0])
			if err != nil {
				return pid, err
			}
			return _pid, nil
		}
	}

	return pid, nil
}

func (s *manager) launchEditor(ctx context.Context, editor editors.Editor, tempDir string) (*exec.Cmd, error) {
	cfgArgs := strings.Split(editor.Args, " ")
	cmdArgs := append(cfgArgs, tempDir)
	cmd := exec.Command(editor.ExecPath, cmdArgs...)
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to run editor: %v", err)
	}

	return cmd, nil
}

func (s *manager) fillWorkspace(ctx context.Context, source Source, tempDir string) error {
	switch source.SourceType {
	case SessionSourceTypeBlank:
		return nil
	case SessionSourceTypeDirectory:
		return s.fillWorkspaceFromDirectory(ctx, source, tempDir)
	case SessionSourceTypeRepository:
		return s.fillWorkspaceFromRepository(ctx, source, tempDir)
	case SessionSourceTypeTemplate:
		return s.fillWorkspaceFromTemplate(ctx, source, tempDir)
	}

	return ErrInvalidSessionSource
}

func (s *manager) fillWorkspaceFromTemplate(ctx context.Context, source Source, tempDir string) error {
	// Check if the specified template exists.
	template, err := s.GetTemplate(ctx, source.Value)
	if err != nil {
		return err
	}

	// Copy the template to the temporary directory
	err = copy.Copy(template.AbsolutePath(ctx), tempDir, copy.Options{
		Skip: func(srcinfo os.FileInfo, src, dest string) (bool, error) {
			// Don't include repository information if the source
			// is a git repository.
			return srcinfo.IsDir() && srcinfo.Name() == ".git", nil
		},
	})
	if err != nil {
		return fmt.Errorf("failed to load template %s: %s", source, err)
	}

	return nil
}

func (s *manager) fillWorkspaceFromRepository(ctx context.Context, source Source, tempDir string) error {
	_, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: source.Value,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	return nil
}

func (s *manager) fillWorkspaceFromDirectory(ctx context.Context, source Source, tempDir string) error {
	absPath, err := filepath.Abs(source.Value)
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
