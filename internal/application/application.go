package application

import (
	"context"
	"errors"
	"os"

	"github.com/fatih/color"

	editor_commands "github.com/letstrygo/letstry/internal/application/commands/editors"
	general_commands "github.com/letstrygo/letstry/internal/application/commands/general"
	hidden_commands "github.com/letstrygo/letstry/internal/application/commands/hidden"
	session_commands "github.com/letstrygo/letstry/internal/application/commands/sessions"
	template_commands "github.com/letstrygo/letstry/internal/application/commands/templates"

	"github.com/letstrygo/letstry/internal/cli"
	"github.com/letstrygo/letstry/internal/environment"
	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/manager"
)

var (
	ErrNoCommandProvided = errors.New("no command provided")
)

type Application struct {
	cli.CliApp

	context context.Context
}

// NewApplication creates a new application instance
func NewApplication(ctx context.Context) *Application {
	// Initialize base line logging, writing to the console.
	logger, err := logging.New(&logging.LoggerConfig{
		LogMode: logging.LogModeConsole,
	})
	if err != nil {
		panic(err)
	}

	// Initialize environment
	ctx = environment.ContextWithEnvironment(ctx)

	// Initialize session manager
	ctx = manager.ContextWithManager(ctx)

	// Initialize logging
	ctx = logging.ContextWithLogger(ctx, logger)

	var commands = []cli.Command{
		session_commands.NewSessionCommand(),
		session_commands.ListSessionsCommand(),
		session_commands.ExportSessionCommand(),
		session_commands.CleanCommand(),
		session_commands.CleanAllCommand(),

		template_commands.ListTemplatesCommand(),
		template_commands.SaveTemplateCommand(),
		template_commands.ImportTemplate(),
		template_commands.DeleteTemplateCommand(),

		editor_commands.ListEditorsCommand(),
		editor_commands.SetEditorCommand(),
		editor_commands.GetEditorCommand(),

		hidden_commands.MonitorCommand(),
		general_commands.VersionCommand(),
	}

	// Initialize the application.
	app := &Application{
		context: ctx,
		CliApp: cli.CliApp{
			Config: cli.CliAppConfig{
				DescriptionMaxWidth: 60,
				HelpCommandSorter:   cli.CommandSorterOrderedAs(commands),
			},

			Name:             cli.MainName(),
			ShortDescription: "a powerful tool for creating temporary workspaces",
			Description:      cli.MainName() + " provides a temporary workspace for you to work in, and then destroys it when you are done.",
		},
	}

	// Add commands
	for _, command := range commands {
		app.RegisterCommand(command)
	}

	// Add help command
	app.RegisterHelpCommand()

	return app
}

// Start starts the application
func (app *Application) Start() {
	// Parse the command line
	invocation, err := app.GetInvocation()
	if err != nil {
		logger, _ := logging.LoggerFromContext(app.GetContext())
		logger.Printf("Error: %s\n", color.RedString(err.Error()))
		os.Exit(1)
	}

	// Configure logging
	err = invocation.UpdateLogger(app)
	if err != nil {
		logger, _ := logging.LoggerFromContext(app.GetContext())
		logger.Printf("Error: %s\n", color.RedString(err.Error()))
		os.Exit(1)
	}
	defer logging.CloseLog(app.GetContext())

	// Run the command
	err = invocation.Execute(app)
	if err != nil {
		logger, _ := logging.LoggerFromContext(app.GetContext())
		logger.Printf("Error: %s\n", color.RedString(err.Error()))
		os.Exit(1)
	}
}

// GetContext returns the application context
func (a *Application) GetContext() context.Context {
	return a.context
}
