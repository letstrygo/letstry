package config

import (
	"fmt"

	"github.com/letstrygo/letstry/internal/config/editors"
)

type Config struct {
	DefaultEditorName editors.EditorName `json:"default_editor"`
	AvailableEditors  []editors.Editor   `json:"editors"`
}

func (cfg Config) GetEditor(name string) (editors.Editor, error) {
	for _, editor := range cfg.AvailableEditors {
		if editor.Name.String() == name {
			return editor, nil
		}
	}

	return editors.Editor{}, fmt.Errorf("editor %s not found", name)
}

func (cfg Config) GetDefaultEditor() (editors.Editor, error) {
	for _, editor := range cfg.AvailableEditors {
		if editor.Name == cfg.DefaultEditorName {
			return editor, nil
		}
	}

	return editors.Editor{}, fmt.Errorf("editor %s not found", cfg.DefaultEditorName)
}
