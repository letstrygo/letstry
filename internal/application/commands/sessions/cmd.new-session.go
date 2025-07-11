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
		ShortDescription: "Create a new session or project",
		Description:      "Create a new session or project using the specified source.",
		Arguments: []cli.Argument{
			{
				Name:        "source",
				Description: "The source to use for the new session or project. Can be a git repository URL, a path to a directory, or the name of a letstry template.\n\nIf source is not provided, the session will be created from a blank source.",
			},
			{
				Name:        "--temp",
				Description: "When set, session will be forcibly stored in a temporary location. This overrides the \"Require Export\" field in your config file.",
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			var source string
			var forceRequireExport bool

			if len(args) > 0 {
				source = args[0]
			}
			if source == "--temp" {
				source = ""
				forceRequireExport = true
			} else if len(args) > 1 {
				forceRequireExport = args[1] == "--temp"
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
				Source:             source,
				ForceRequireExport: forceRequireExport,
			})
			if err != nil {
				return err
			}

			if session != nil {
				logger.Printf("session created: %s\n", session.String())
			} else {
				logger.Printf("project created")
			}

			return nil
		},
	}
}
