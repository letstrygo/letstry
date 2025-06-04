package cli

import (
	"context"
	"os"
)

const HelpCommandName string = "help"

func helpCommand(app CliApp) Command {
	return Command{
		Name:             HelpCommandName,
		Description:      "Display help information for the application or a specific command.",
		ShortDescription: "Display help information",
		Aliases:          []string{"-h", "--help"},
		Arguments: []Argument{
			{
				Name:        "command",
				Description: "The command to display help information for.",
				Required:    false,
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			var inputCmd string

			if len(args) > 0 {
				inputCmd = args[0]
			}

			if inputCmd != "" {
				if !app.IsCommand(inputCmd) {
					return ErrUnknownCommand
				}

				cmd, err := app.Command(inputCmd)
				if err != nil {
					return err
				}
				return cmd.Write(app, os.Stdout)
			}

			return app.Write(os.Stdout)
		},
	}
}
