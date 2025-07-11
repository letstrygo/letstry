package general

import (
	"context"

	"github.com/letstrygo/letstry/internal/application/commands"
	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/config"
	"github.com/letstrygo/letstry/internal/storage"
)

func PathCommand() cli.Command {
	return cli.Command{
		Name:             commands.CommandPath.String(),
		ShortDescription: "Displays paths relevant to LetsTry",
		Description:      "Useful for locating your LetsTry config file, your LetsTry session cache file, or the directory in which LetsTry stores new sessions and projects.",
		Arguments: []cli.Argument{
			{
				Name:        "file",
				Description: "Can be one of config, projects or sessions. Defaults to \"config\".",
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}

			var opt string
			if len(args) > 0 {
				opt = args[0]
			}

			switch opt {
			case "sessions":
				store := storage.GetStorage()
				sessPath := store.GetAbsolutePath("sessions.json")
				println(sessPath)
				return nil
			case "projects":
				path := cfg.LTPath
				if cfg.RequireExport {
					println("require-export enabled, all projects will be stored in a temporary directory and deleted once their corresponding session is terminated.")
					return nil
				}

				if path == "" {
					println("project-path: custom projects_path not set, all projects will be stored in a temporary directory and deleted once their corresponding session is terminated.")
					return nil
				}

				println(cfg.LTPath)
				return nil
			default:
				println(cfg.Path())
				return nil
			}
		},
	}
}
