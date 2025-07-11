package config

import (
	"fmt"

	"github.com/letstrygo/letstry/internal/config/editors"
)

type Config struct {
	path string

	// The path at which LT projects will be placed. Defaults to your systems
	// temporary directory.
	LTPath string `json:"projects_path"`
	// When enabled, projects created in letstry will be deleted after the
	// lt session is closed (i.e. when the editor is closed).
	//
	// You can override this by passing `--temp` when creating the LetsTry session.
	RequireExport bool `json:"require_export"`
	// The name of the default editor to use for new sessions. (Default: vscode)
	DefaultEditorName editors.EditorName `json:"default_editor"`
	// Editors available for use within LetsTry. (Default: vscode)
	AvailableEditors []editors.Editor `json:"editors"`
}

func (cfg Config) Path() string {
	return cfg.path
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
