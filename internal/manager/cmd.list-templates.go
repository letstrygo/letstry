package manager

import "context"

func (s *manager) ListTemplates(ctx context.Context) ([]Template, error) {
	if !s.storage.DirectoryExists("templates") {
		return []Template{}, nil
	}

	templates, err := s.storage.ListDirectories("templates")
	if err != nil {
		return nil, err
	}

	var result []Template
	for _, t := range templates {
		result = append(result, Template(t))
	}

	return result, nil
}
