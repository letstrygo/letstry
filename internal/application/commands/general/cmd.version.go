package general

import (
	"context"
	"debug/buildinfo"
	"fmt"
	"os"
	"runtime"

	"github.com/letstrygo/letstry/internal/application/commands"
	"github.com/letstrygo/letstry/internal/cli"
)

func VersionCommand() cli.Command {
	return cli.Command{
		Name:             commands.CommandVersion.String(),
		Aliases:          []string{"-v", "--version"},
		ShortDescription: "Display the version of " + cli.MainName(),
		Description:      "This command will display the version of " + cli.MainName() + ".",
		Executor: func(ctx context.Context, args []string) error {
			exe, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %v", err)
			}

			info, err := buildinfo.ReadFile(exe)
			if err != nil {
				return fmt.Errorf("failed to read build info: %v", err)
			}

			fmt.Println(cli.MainName(), "version", info.Main.Version, fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))

			return nil
		},
	}
}
