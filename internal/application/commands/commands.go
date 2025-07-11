package commands

type CommandName string

func (c CommandName) String() string {
	return string(c)
}

const (
	CommandVersion        CommandName = "version"
	CommandPath           CommandName = "path"
	CommandMonitor        CommandName = "monitor"
	CommandClean          CommandName = "clean"
	CommandPruneSessions  CommandName = "prune"
	CommandNewSession     CommandName = "new"
	CommandListSessions   CommandName = "list"
	CommandListTemplates  CommandName = "templates"
	CommandListEditors    CommandName = "editors"
	CommandGetEditor      CommandName = "get-editor"
	CommandSetEditor      CommandName = "set-editor"
	CommandDeleteTemplate CommandName = "delete"
	CommandSaveTemplate   CommandName = "save"
	CommandUpdateTemplate CommandName = "update"
	CommandExportSession  CommandName = "export"
	CommandShow           CommandName = "show"
)
