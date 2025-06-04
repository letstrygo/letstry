package hidden

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/letstrygo/letstry/internal/application/commands"
	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/config/editors"
	"github.com/letstrygo/letstry/internal/manager"
)

var (
	ErrMissingArgumentDelay        = errors.New("monitor: missing required argument 'delay'")
	ErrMissingArgumentLocation     = errors.New("monitor: missing required argument 'location'")
	ErrMissingArgumentPID          = errors.New("monitor: missing required argument 'pid'")
	ErrMissingArgumentTrackingType = errors.New("monitor: missing required argument 'tracking-type'")
)

func MonitorCommand() cli.Command {
	return cli.Command{
		Name:             commands.CommandMonitor.String(),
		ShortDescription: "Monitor a session",
		Description:      "This command is used to monitor a session. It is not intended to be run directly by the user.",
		LogToFile:        true,
		Arguments: []cli.Argument{
			{
				Name:        "delay",
				Description: "The delay before initiating the monitor, formatted as a duration string. For example, '5s' for 5 seconds.",
				Required:    true,
			},
			{
				Name:        "location",
				Description: "The location of the session.",
				Required:    true,
			},
			{
				Name:        "pid",
				Description: "The process ID to monitor.",
				Required:    true,
			},
			{
				Name:        "tracking-type",
				Description: "The type of tracking data required for the monitor to run. Can be one of the following: 'file_access', 'process'.",
				Required:    true,
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			if len(args) < 1 {
				return ErrMissingArgumentDelay
			}

			if len(args) < 2 {
				return ErrMissingArgumentLocation
			}

			if len(args) < 3 {
				return ErrMissingArgumentPID
			}

			if len(args) < 4 {
				return ErrMissingArgumentTrackingType
			}

			delay, err := time.ParseDuration(args[0])
			if err != nil {
				return err
			}

			location := args[1]
			pid, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			mgr, err := manager.GetManager(ctx)
			if err != nil {
				return err
			}

			trackingType, err := editors.GetTrackingType(args[3])
			if err != nil {
				return err
			}

			return mgr.MonitorSession(ctx, manager.MonitorSessionArguments{
				Delay:        delay,
				TrackingType: trackingType,
				PID:          pid,
				Location:     location,
			})
		},
	}
}
