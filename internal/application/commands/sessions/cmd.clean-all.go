package sessions

import (
	"context"

	"github.com/letstrygo/letstry/internal/application/commands"
	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/manager"
)

func CleanAllCommand() cli.Command {
	return cli.Command{
		Name:             commands.CommandCleanAll.String(),
		ShortDescription: "Clean any dangling sessions",
		Description:      "This command will clean any dangling sessions that are no longer in use. This command is useful for cleaning up any sessions that were not properly closed.",
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

			inactiveSessions := []manager.Session{}
			for _, session := range sessions {
				if !session.IsActive() {
					inactiveSessions = append(inactiveSessions, session)
				}
			}

			if len(inactiveSessions) < 1 {
				logger.Printf("%s: no inactive sessions to clean\n", commands.CommandCleanAll)
				return nil
			}

			for _, session := range inactiveSessions {
				err := mgr.CleanSession(ctx, manager.CleanSessionArguments{
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
