package commands

type CommandName string

func (c CommandName) String() string {
	return string(c)
}

const (
	CommandVersion        CommandName = "version"
	CommandMonitor        CommandName = "monitor"
	CommandClean          CommandName = "clean"
	CommandCleanAll       CommandName = "clean-all"
	CommandNewSession     CommandName = "new"
	CommandListSessions   CommandName = "list"
	CommandListTemplates  CommandName = "templates"
	CommandListEditors    CommandName = "editors"
	CommandGetEditor      CommandName = "get-editor"
	CommandSetEditor      CommandName = "set-editor"
	CommandDeleteTemplate CommandName = "delete-template"
	CommandSaveTemplate   CommandName = "save"
	CommandExportSession  CommandName = "export"
)
