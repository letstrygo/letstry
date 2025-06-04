package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"text/template"
)

type CliAppConfig struct {
	DescriptionMaxWidth int
	HelpCommandSorter   CommandSorter
}

type CliApp struct {
	Config           CliAppConfig
	Name             string
	ShortDescription string
	Description      string

	commands map[string]Command
}

func (app CliApp) Write(out io.Writer) error {
	return template.Must(
		template.
			New("help").
			Funcs(defaultTemplateFuncs(app)).
			Parse(applicationHelpTemplate),
	).Execute(out, app)
}

func (app CliApp) IsCommand(name string) bool {
	if app.commands == nil {
		return false
	}

	for _, command := range app.commands {
		if name == command.Name || slices.Contains(command.Aliases, name) {
			return true
		}
	}

	return false
}

func (app CliApp) Command(name string) (Command, error) {
	if app.commands == nil {
		return Command{}, ErrUnknownCommand
	}

	for _, command := range app.commands {
		if name == command.Name || slices.Contains(command.Aliases, name) {
			name = command.Name
		}
	}

	if _, ok := app.commands[name]; !ok {
		return Command{}, ErrUnknownCommand
	}

	return app.commands[name], nil
}

func (app *CliApp) RegisterCommand(command Command) error {
	err := app.registerCommand(command)
	if err != nil {
		return err
	}

	// Re-register the help command if it's been registered.
	if _, helpCommandRegistered := app.commands[HelpCommandName]; helpCommandRegistered {
		app.RegisterHelpCommand()
	}

	return nil
}

func (app *CliApp) registerCommand(command Command) error {
	if app.commands == nil {
		app.commands = make(map[string]Command)
	}

	if _, ok := app.commands[command.Name]; ok {
		return fmt.Errorf("command %s already registered", command.Name)
	}

	app.commands[command.Name] = command

	return nil
}

func (app *CliApp) RegisterHelpCommand() {
	if app.commands != nil {
		delete(app.commands, HelpCommandName)
	}

	app.registerCommand(helpCommand(*app))
}

func MainName() string {
	cmd := filepath.Base(os.Args[0])

	if runtime.GOOS == "windows" {
		cmd = strings.TrimSuffix(cmd, ".exe")
	}

	return cmd
}
