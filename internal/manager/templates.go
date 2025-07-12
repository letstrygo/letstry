package manager

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/letstrygo/templates"
)

type Template struct {
	*templates.Template
}

func (t Template) String() string {
	return string(t.Name)
}

func (t Template) FormattedString(ctx context.Context) string {
	name := color.YellowString(t.String())

	typeStr := string(t.Type)
	switch t.Type {
	case templates.TemplateTypeGitRepository:
		typeStr = color.MagentaString(string(t.Type))
	case templates.TemplateTypeLocal:
		typeStr = color.HiBlueString(string(t.Type))
	}

	locationStr := color.BlueString("%5v", t.Source)

	// updated = color.BlueString("(%s)", t..Format("2006-01-02 15:04:05"))
	// TODO: Add UpdatedBy to templates database

	detailsStr := fmt.Sprintf("[%v=%v, %v=%v]", color.HiBlackString("type"), typeStr, color.HiBlackString("location"), locationStr)
	return fmt.Sprintf("%s %v", name, detailsStr)
}

func (s *manager) GetTemplate(ctx context.Context, name string) (Template, error) {
	tmpl, err := s.repository.GetTemplateByName(name)
	if err != nil {
		return Template{}, err
	}

	return Template{tmpl}, nil
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
