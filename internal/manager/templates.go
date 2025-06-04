package manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type Template string

func (t Template) String() string {
	return string(t)
}

func (t Template) FormattedString(ctx context.Context) string {
	name := color.YellowString(t.String())

	var updated string

	absolutePath := t.AbsolutePath(ctx)
	// Get the last updated date for the directory
	stat, err := os.Stat(absolutePath)
	if err != nil {
		updated = color.RedString("unknown")
	} else {
		updated = color.BlueString("(%s)", stat.ModTime().Format("2006-01-02 15:04:05"))
	}

	return fmt.Sprintf("name=%s, updated=%s", name, updated)
}

func (t Template) AbsolutePath(ctx context.Context) string {
	sessionMgr, err := GetManager(ctx)
	if err != nil {
		panic(err)
	}

	return sessionMgr.storage.GetAbsolutePath(t.StoragePath())
}

func (t Template) StoragePath() string {
	return filepath.Join("templates", t.String())
}

func (s *manager) GetTemplate(ctx context.Context, name string) (Template, error) {
	template := Template(name)

	if !s.storage.DirectoryExists(template.StoragePath()) {
		return "", fmt.Errorf("template with name %s does not exist", name)
	}

	return template, nil
}

func (s *manager) createTemplatesDirectoryIfNotExists() error {
	if !s.storage.DirectoryExists("templates") {
		err := s.storage.CreateDirectory("templates")
		if err != nil {
			return err
		}
	}

	return nil
}
