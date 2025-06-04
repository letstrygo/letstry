package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/process"

	"github.com/letstrygo/letstry/internal/config/editors"
	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/util/access"
	"github.com/letstrygo/letstry/internal/util/identifier"
)

var (
	ErrUnknownTrackingType = errors.New("unknown tracking type")
)

type MonitorSessionArguments struct {
	Delay        time.Duration
	TrackingType editors.TrackingType
	PID          int
	Location     string
}

func (s *manager) MonitorSession(ctx context.Context, args MonitorSessionArguments) error {
	// delay the start of the monitoring
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	logger.Printf("delaying monitoring for %v\n", args.Delay)
	time.Sleep(args.Delay)

	session, err := s.GetSessionForPath(ctx, args.Location)
	if err != nil {
		return err
	}
	logger.Printf("monitoring session: %s\n", session.ID)

	logger = logger.ChildLogger(fmt.Sprintf("sess-%s", session.ID))

	handler := func() error {
		switch session.Editor.TrackingType {
		case editors.TrackingTypeFileAccess:
			logger.Printf("cleaning up session: %s (directory no longer being accessed)\n", session.ID)
		case editors.TrackingTypeProcess:
			logger.Printf("cleaning up session: %s (process no longer running)\n", session.ID)
		}

		err = s.removeSession(ctx, session.ID)
		if err != nil {
			return err
		}

		return nil
	}

	switch args.TrackingType {
	case editors.TrackingTypeProcess:
		logger.Printf("using tracking type: %v\n", editors.TrackingTypeProcess)
		return s.monitorProcess(args.PID, handler)
	case editors.TrackingTypeFileAccess:
		logger.Printf("using tracking type: %v\n", editors.TrackingTypeFileAccess)
		_, err := os.Stat(args.Location)
		if err != nil {
			return err
		}
		return s.monitorDirectoryAccessible(args.Location, handler)
	}

	return ErrUnknownTrackingType
}

func (s *manager) monitorProcess(pid int, callback func() error) error {
	for {
		_, err := process.NewProcess(int32(pid))
		if err != nil {
			return callback()
		}

		time.Sleep(1 * time.Second) // Check every second
	}
}

func (s *manager) monitorDirectoryAccessible(path string, callback func() error) error {
	for {
		if !access.IsPathUse(path) {
			return callback()
		}

		time.Sleep(1 * time.Second) // Check every second
	}
}

func (s *manager) removeSession(ctx context.Context, id identifier.ID) error {
	sessions, err := s.ListSessions(ctx)
	if err != nil {
		return err
	}

	for i, session := range sessions {
		if session.ID == id {
			// Remove the session
			sessions = append(sessions[:i], sessions[i+1:]...)

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
				return fmt.Errorf("failed to truncate sessions file: %v", err)
			}

			_, err = file.Write(data)
			if err != nil {
				return fmt.Errorf("failed to write sessions: %v", err)
			}

			err = file.Sync()
			if err != nil {
				return fmt.Errorf("failed to sync sessions file: %v", err)
			}

			// Give the process manager time to settle
			time.Sleep(1 * time.Second)

			// Remove the temporary directory
			err = os.RemoveAll(session.Location)
			if err != nil {
				return fmt.Errorf("failed to remove temporary directory: %v", err)
			}

			return nil
		}
	}

	return fmt.Errorf("session with id %s not found", id)
}
