package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/letstrygo/letstry/internal/logging"
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

	sourceType, err := s.GetSessionSourceType(ctx, args.RepositoryUrl)
	if err != nil {
		return zeroValue, err
	}

	if sourceType != SessionSourceTypeRepository {
		return zeroValue, ErrInvalidSourceType
	}

	template := Template(args.TemplateName)

	if s.storage.DirectoryExists(template.StoragePath()) {
		return zeroValue, fmt.Errorf("template already exists: %s", template.String())
	}

	err = s.createTemplatesDirectoryIfNotExists()
	if err != nil {
		return "", err
	}

	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return "", err
	}

	logger.Printf("cloning repository %s\n", args.RepositoryUrl)
	_, err = git.PlainClone(template.AbsolutePath(ctx), false, &git.CloneOptions{
		URL: args.RepositoryUrl,
	})
	if err != nil {
		return zeroValue, fmt.Errorf("failed to clone repository: %v", err)
	}

	logger.Printf("imported template: %s\n", template.FormattedString(ctx))
	return template, nil
}
