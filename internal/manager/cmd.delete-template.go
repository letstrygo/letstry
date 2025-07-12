package manager

import (
	"context"
	"path/filepath"

	"github.com/letstrygo/letstry/internal/logging"
)

func (s *manager) DeleteTemplate(ctx context.Context, template Template) error {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	storagePath := filepath.Join("templates", template.Name)
	if s.storage.DirectoryExists(storagePath) {
		s.storage.DeleteDirectory(storagePath)
	}

	err = s.repository.DeleteTemplate(template.Name)
	if err != nil {
		return err
	}

	logger.Printf("deleted template: %s\n", template.String())
	return nil
}
