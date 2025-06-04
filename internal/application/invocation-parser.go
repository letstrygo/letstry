package application

import (
	"os"

	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/logging"
)

type ApplicationInvocation struct {
	Command   cli.Command
	Arguments []string
}

func (i ApplicationInvocation) Execute(app *Application) error {
	return i.Command.Execute(app.GetContext(), i.Arguments)
}

func (i ApplicationInvocation) UpdateLogger(app *Application) error {
	if i.Command.LogToFile {
		logger, err := logging.New(&logging.LoggerConfig{
			LogMode: logging.LogModeFile,
		})
		if err != nil {
			return err
		}

		app.context = logging.ContextWithLogger(app.GetContext(), logger)
	}

	return nil
}

func (a *Application) GetInvocation() (*ApplicationInvocation, error) {
	var err error
	var command cli.Command
	var cmdArgs []string = []string{string(cli.HelpCommandName)}

	args := os.Args

	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	command, err = a.Command(cmdArgs[0])
	if err != nil {
		return nil, err
	}

	return &ApplicationInvocation{
		Command:   command,
		Arguments: cmdArgs[1:],
	}, nil
}
