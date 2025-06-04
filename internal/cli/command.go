package cli

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/letstrygo/letstry/internal/manager"
)

type CommandExecutor func(context.Context, []string) error

type Command struct {
	Name                 string
	ShortDescription     string
	Description          string
	Aliases              []string
	Arguments            []Argument
	Executor             CommandExecutor
	LogToFile            bool
	MustBeRunFromSession bool
}

func (command Command) Execute(ctx context.Context, args []string) error {
	if command.MustBeRunFromSession {
		mgr, err := manager.GetManager(ctx)
		if err != nil {
			return err
		}

		_, err = mgr.GetCurrentSession(ctx)
		if err != nil {
			return err
		}
	}

	return command.Executor(ctx, args)
}

func (command Command) Write(app CliApp, out io.Writer) error {
	return template.Must(
		template.
			New(fmt.Sprintf("help.%s", command.Name)).
			Funcs(defaultTemplateFuncs(app)).
			Parse(commandHelpTemplate),
	).Execute(out, command)
}
