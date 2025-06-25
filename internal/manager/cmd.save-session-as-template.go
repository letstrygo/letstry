package manager

import (
	"context"
	"errors"

	"github.com/letstrygo/letstry/internal/logging"
	"github.com/otiai10/copy"
)

var (
	ErrMissingTemplateName = errors.New("missing template name")
)

type SaveSessionAsTemplateArguments struct {
	TemplateName string `json:"template_name"`
}

func (s *manager) SaveSessionAsTemplate(ctx context.Context, arg SaveSessionAsTemplateArguments) (Template, error) {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return "", err
	}

	session, err := s.GetCurrentSession(ctx)
	if err != nil {
		return "", err
	}

	err = s.createTemplatesDirectoryIfNotExists()
	if err != nil {
		return "", err
	}

	var template Template

	if arg.TemplateName != "" {
		template = Template(arg.TemplateName)
	} else {
		if session.Source.SourceType == SessionSourceTypeTemplate {
			template = Template(session.Source.Value)
		}
	}

	if template == "" {
		return "", ErrMissingTemplateName
	}

	if s.storage.DirectoryExists(template.StoragePath()) {
		logger.Printf("template already exists, deleting template %s\n", template.String())
		err = s.storage.DeleteDirectory(template.StoragePath())
		if err != nil {
			return "", err
		}
	}

	logger.Printf("creating template %s from session %s\n", template.String(), session.ID.FormattedString())

	err = s.storage.CreateDirectory(template.StoragePath())
	if err != nil {
		return "", err
	}

	err = copy.Copy(session.Location, template.AbsolutePath(ctx))
	if err != nil {
		return "", err
	}

	return template, nil
}
