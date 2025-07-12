package manager

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/templates"
)

var (
	ErrInvalidSourceType = errors.New("invalid source type, must be repository")
)

type ImportTemplateArguments struct {
	TemplateName  string
	RepositoryUrl string
}

func (s *manager) ImportTemplate(ctx context.Context, args ImportTemplateArguments) (Template, error) {
	var zeroValue Template

	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return zeroValue, err
	}

	sourceType, err := s.GetSessionSourceType(ctx, args.RepositoryUrl)
	if err != nil {
		return zeroValue, err
	}

	if sourceType != SessionSourceTypeRepository {
		return zeroValue, ErrInvalidSourceType
	}

	t, err := s.repository.GetTemplateByName(args.TemplateName)
	if err != nil {
		if !errors.Is(err, templates.ErrTemplateNotFound) {
			return zeroValue, err
		}
	}

	if t != nil {
		logger.Printf("a template with that name already exists\n")
		return Template{}, templates.ErrTemplateExists
	}

	storagePath := filepath.Join("templates", args.TemplateName)
	absoluteStoragePath := s.storage.GetAbsolutePath(storagePath)

	if s.storage.DirectoryExists(absoluteStoragePath) {
		logger.Printf("a local template directory already exists with that name but is not tracked by the database\n")
		logger.Printf("please delete this directory to continue\n")
		logger.Printf("%v\n", absoluteStoragePath)

		return Template{}, ErrUntrackedTemplate
	}

	err = s.repository.CreateTemplate(templates.CreateTemplate{
		Name:       args.TemplateName,
		Source:     args.RepositoryUrl,
		Type:       templates.TemplateTypeGitRepository,
		IsOfficial: false,
	})
	if err != nil {
		return zeroValue, err
	}

	template, err := s.GetTemplate(ctx, args.TemplateName)
	if err != nil {
		return zeroValue, err
	}

	logger.Printf("imported template: %s\n", template.FormattedString(ctx))

	return template, nil
}
