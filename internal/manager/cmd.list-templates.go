package manager

import (
	"context"

	"github.com/letstrygo/templates"
	"github.com/samber/lo"
)

func (s *manager) ListTemplates(ctx context.Context) ([]Template, error) {
	tmpls, err := s.repository.ListTemplates(templates.ListTemplates{})
	if err != nil {
		return nil, err
	}

	return lo.Map(tmpls, func(t templates.Template, _ int) Template {
		return Template{&t}
	}), nil
}
