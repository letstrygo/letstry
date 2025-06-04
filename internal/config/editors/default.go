package editors

func DefaultEditors() []Editor {
	allEditors := []Editor{
		VSCodeEditor(),
		NotepadPlusPlusEditor(),
	}

	result := []Editor{}

	for _, editor := range allEditors {
		// Check if the editor is installed
		if editor.IsInstalled() {
			result = append(result, editor)
		}
	}

	return result
}
