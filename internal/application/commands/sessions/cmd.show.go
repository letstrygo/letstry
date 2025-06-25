package sessions

import (
	"context"

	"github.com/letstrygo/letstry/internal/application/commands"
	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/manager"
	"github.com/letstrygo/letstry/internal/util/identifier"
)

func ShowCommand() cli.Command {
	return cli.Command{
		Name:                 commands.CommandShow.String(),
		ShortDescription:     "Show details about the session",
		Description:          "This command must be run from within a session. It will retrieve information from the current session. If you are not currently in a session, you can pass a session ID to this comand.",
		MustBeRunFromSession: false,
		Arguments: []cli.Argument{
			{
				Name:        "session-id",
				Description: "The session ID to show information for",
				Required:    false,
			},
			{
				Name:        "output",
				Description: "Either 'full' (default), 'path', 'pid', 'editor' or 'json'",
				Required:    false,
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			var (
				displayType manager.SessionDisplayType = manager.SessionDisplayTypeFull
				sessionID   *identifier.ID
			)

			if len(args) > 0 {
				// If first arg is a display type, don't treat it as an ID
				if dt, err := manager.ParseSessionDisplayType(args[0]); err == nil {
					displayType = dt
				} else {
					sessionID = identifier.ParseIDPtr(args[0])
				}
			}

			if len(args) > 1 {
				dt, err := manager.ParseSessionDisplayType(args[1])
				if err != nil {
					return err
				}

				displayType = dt
			}

			mgr, err := manager.GetManager(ctx)
			if err != nil {
				return err
			}

			return mgr.ShowSession(ctx, manager.ShowSessionArguments{
				SessionID:  sessionID,
				OutputType: displayType,
			})
		},
	}
}
