package sessions

import (
	"context"

	"github.com/letstrygo/letstry/internal/application/commands"
	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/manager"
	"github.com/letstrygo/letstry/internal/util/identifier"
)

func PruneSessionsCommand() cli.Command {
	return cli.Command{
		Name: commands.CommandPruneSessions.String(),
		Aliases: []string{
			// For backwards compatibility
			"clean-all",
			"clean",
		},
		ShortDescription: "Prune inactive sessions",
		Description:      "This command will prune any inactive sessions that are no longer in use. This command is useful for cleaning up any sessions that were not properly closed.",
		Arguments: []cli.Argument{
			{
				Name:        "session-id",
				Description: "The session to prune. (Defaults to all inactive sessions)",
				Required:    false,
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			mgr, err := manager.GetManager(ctx)
			if err != nil {
				return err
			}

			logger, err := logging.LoggerFromContext(ctx)
			if err != nil {
				return err
			}

			sessions, err := mgr.ListSessions(ctx)
			if err != nil {
				return err
			}

			if len(args) > 0 {
				id := identifier.ID(args[0])
				sess, err := mgr.GetSession(ctx, id)
				if err != nil {
					return err
				}
				sessions = []manager.Session{sess}
			}

			inactiveSessions := []manager.Session{}
			for _, session := range sessions {
				if !session.IsActive() {
					inactiveSessions = append(inactiveSessions, session)
				}
			}

			if len(inactiveSessions) < 1 {
				logger.Printf("%s: no inactive sessions to prune\n", commands.CommandPruneSessions)
				return nil
			}

			for _, session := range inactiveSessions {
				err := mgr.PruneSession(ctx, manager.PruneSessionArguments{
					SessionID: session.ID,
				})
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
