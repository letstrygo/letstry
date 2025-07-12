package manager

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/templates"
	"github.com/otiai10/copy"
)

var (
	ErrMissingTemplateName        = errors.New("missing template name")
	ErrCannotOverwriteGitTemplate = errors.New("cannot overwrite git template")
	ErrUntrackedTemplate          = errors.New("untracked template")
)

type SaveSessionAsTemplateArguments struct {
	TemplateName string `json:"template_name"`
}

func (s *manager) SaveSessionAsTemplate(ctx context.Context, arg SaveSessionAsTemplateArguments) (Template, error) {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return Template{}, err
	}

	session, err := s.GetCurrentSession(ctx)
	if err != nil {
		return Template{}, err
	}

	err = s.createTemplatesDirectoryIfNotExists()
	if err != nil {
		return Template{}, err
	}

	var templateName string

	if arg.TemplateName != "" {
		templateName = arg.TemplateName
	} else {
		if session.Source.SourceType == SessionSourceTypeTemplate {
			templateName = session.Source.Value
		}
	}

	if templateName == "" {
		return Template{}, ErrMissingTemplateName
	}

	template, err := s.GetTemplate(ctx, templateName)
	if err != nil {
		if !errors.Is(err, templates.ErrTemplateNotFound) {
			return Template{}, err
		}
	}

	var (
		storagePath         = filepath.Join("templates", templateName)
		absoluteStoragePath = s.storage.GetAbsolutePath(storagePath)
	)

	if template.Template != nil {
		if template.Type == templates.TemplateTypeGitRepository {
			logger.Printf("%v is a git template, these cannot be overwritten. please use another template name", templateName)
			return Template{}, ErrCannotOverwriteGitTemplate
		}

		logger.Printf("template already exists, deleting template %s\n", template.String())
		err = s.DeleteTemplate(ctx, template)
		if err != nil {
			return Template{}, err
		}
	}

	if s.storage.DirectoryExists(storagePath) {
		logger.Printf("a local template directory already exists with that name but is not tracked by the database\n")
		logger.Printf("please delete this directory to continue\n")
		logger.Printf("%v\n", absoluteStoragePath)

		return Template{}, ErrUntrackedTemplate
	}

	logger.Printf("creating template %s from session %s\n", templateName, session.ID.FormattedString())

	err = s.storage.CreateDirectory(storagePath)
	if err != nil {
		return Template{}, err
	}

	err = copy.Copy(session.Location, absoluteStoragePath)
	if err != nil {
		return Template{}, err
	}

	err = s.repository.CreateTemplate(templates.CreateTemplate{
		Name:       templateName,
		Source:     absoluteStoragePath,
		Type:       templates.TemplateTypeLocal,
		IsOfficial: false,
	})
	if err != nil {
		return Template{}, err
	}

	return Template{}, nil
}
