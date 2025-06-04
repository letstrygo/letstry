package manager

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

var (
	ErrInvalidSessionSource = errors.New("invalid session source")
)

type SessionSourceType string

func (s SessionSourceType) String() string {
	return string(s)
}

func (s SessionSourceType) FormattedString() string {
	return color.HiYellowString(s.String())
}

const (
	SessionSourceTypeRepository SessionSourceType = "repository"
	SessionSourceTypeDirectory  SessionSourceType = "directory"
	SessionSourceTypeTemplate   SessionSourceType = "template"
	SessionSourceTypeBlank      SessionSourceType = "blank"
)

// GetSessionSourceType returns the type of session source for the given value.
func (s *manager) GetSessionSourceType(ctx context.Context, value string) (SessionSourceType, error) {
	// Check for blank value.
	if value == "" {
		return SessionSourceTypeBlank, nil
	}

	// Check for template.
	_, err := s.GetTemplate(ctx, value)
	if err == nil {
		return SessionSourceTypeTemplate, nil
	}

	// Check if directory exists and is a directory.
	if absPath, err := filepath.Abs(value); err == nil {
		stat, err := os.Stat(absPath)
		if err == nil && stat.IsDir() {
			return SessionSourceTypeDirectory, nil
		}
	}

	// Check for repository.
	_, err = git.NewRemote(nil, &config.RemoteConfig{
		URLs: []string{value},
	}).List(&git.ListOptions{})
	if err == nil {
		return SessionSourceTypeRepository, nil
	}

	return "", ErrInvalidSessionSource
}
