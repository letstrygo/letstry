package manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/letstrygo/letstry/internal/logging"
	"github.com/otiai10/copy"
)

type ExportSessionArguments struct {
	Path string
}

func (s *manager) ExportSession(ctx context.Context, arg ExportSessionArguments) error {
	session, err := s.GetCurrentSession(ctx)
	if err != nil {
		return err
	}

	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(arg.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(absPath); err == nil {
		return fmt.Errorf("path %s already exists", absPath)
	}

	logger.Printf("exporting session %s to %s\n", session.FormattedID(), absPath)
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy the session to the export path
	if err := copy.Copy(session.Location, absPath); err != nil {
		return fmt.Errorf("failed to copy session: %w", err)
	}
	logger.Printf("session exported successfully\n")

	return nil
}
