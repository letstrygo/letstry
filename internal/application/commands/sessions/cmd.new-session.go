package sessions

import (
	"context"

	"github.com/letstrygo/letstry/internal/application/commands"
	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/manager"
)

// NewSessionCommand returns a new command for creating a new session.
func NewSessionCommand() cli.Command {
	return cli.Command{
		Name:             commands.CommandNewSession.String(),
		ShortDescription: "Create a new session",
		Description:      "Create a new session using the specified source.",
		Arguments: []cli.Argument{
			{
				Name:        "source",
				Description: "The source to use for the new session. Can be a git repository URL, a path to a directory, or the name of a letstry template.\n\nIf source is not provided, the session will be created from a blank source.",
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			var source string

			if len(args) >= 1 {
				source = args[0]
			}

			mgr, err := manager.GetManager(ctx)
			if err != nil {
				return err
			}

			logger, err := logging.LoggerFromContext(ctx)
			if err != nil {
				return err
			}

			session, err := mgr.CreateSession(ctx, manager.CreateSessionArguments{
				Source: source,
			})
			if err != nil {
				return err
			}

			logger.Printf("session created: %s\n", session.String())

			return nil
		},
	}
}
