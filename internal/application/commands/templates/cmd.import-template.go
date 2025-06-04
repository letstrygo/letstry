package templates

import (
	"context"
	"errors"

	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/manager"
)

var (
	ErrMissingRepository = errors.New("missing repository")
)

func ImportTemplate() cli.Command {
	return cli.Command{
		Name:             "import",
		ShortDescription: "Import a template from a git repository",
		Description:      "This command allows you to import a template from a git repository.",
		Arguments: []cli.Argument{
			{
				Name:        "template-name",
				Description: "The name to use for the template.",
				Required:    true,
			},
			{
				Name:        "repository",
				Description: "The git repository to import the template from.",
				Required:    true,
			},
		},
		Executor: func(ctx context.Context, args []string) error {
			mgr, err := manager.GetManager(ctx)
			if err != nil {
				return err
			}

			if len(args) < 1 {
				return ErrMissingTemplateName
			}

			templateName := args[0]

			if len(args) < 2 {
				return ErrMissingRepository
			}

			repository := args[1]

			_, err = mgr.ImportTemplate(ctx, manager.ImportTemplateArguments{
				TemplateName:  templateName,
				RepositoryUrl: repository,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}
}
