package templates

import (
	"context"

	"github.com/letstrygo/letstry/internal/application/commands"
	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/manager"
)

func SaveTemplateCommand() cli.Command {
	return cli.Command{
		Name:                 commands.CommandSaveTemplate.String(),
		ShortDescription:     "Saves the current session as a template",
		Description:          "This command must be run from within a session. It will save the current session as a template with the specified name. If no name is provided, and the session was created from a template, the template's name will be used.",
		MustBeRunFromSession: true,
		Arguments: []cli.Argument{
			{
				Name:        "template-name",
				Description: "The name to use for the template. If not provided, and the session was created from a template, the template's name will be used.",
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			var templateName string

			if len(args) >= 1 {
				templateName = args[0]
			}

			mgr, err := manager.GetManager(ctx)
			if err != nil {
				return err
			}

			_, err = mgr.SaveSessionAsTemplate(ctx, manager.SaveSessionAsTemplateArguments{
				TemplateName: templateName,
			})

			if err != nil {
				return err
			}

			return nil
		},
	}
}
