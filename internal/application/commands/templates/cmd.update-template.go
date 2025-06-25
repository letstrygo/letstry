package templates

import (
	"context"

	"github.com/letstrygo/letstry/internal/application/commands"
	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/manager"
)

func UpdateTemplateCommand() cli.Command {
	return cli.Command{
		Name:                 commands.CommandUpdateTemplate.String(),
		ShortDescription:     "Updates the specified template (if it's a git repository)",
		Description:          "If the specified template is a git repository it will be updated with the latest remote changes",
		MustBeRunFromSession: false,
		Arguments: []cli.Argument{
			{
				Name:        "template-name",
				Description: "The name of the template to update.",
				Required:    true,
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			templateName := args[0]

			mgr, err := manager.GetManager(ctx)
			if err != nil {
				return err
			}

			err = mgr.UpdateTemplate(ctx, manager.UpdateTemplateArguments{
				TemplateName: templateName,
			})
			if err != nil {
				return err
			}

			t, err := mgr.GetTemplate(ctx, templateName)
			if err != nil {
				return err
			}

			log, err := logging.LoggerFromContext(ctx)
			if err != nil {
				return err
			}

			log.Printf("updated template %v", t.FormattedString(ctx))

			return nil
		},
	}
}
