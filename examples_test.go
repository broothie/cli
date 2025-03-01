package cli_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/broothie/cli"
)

func Example_basic_usage() {
	// Create and run a command called "fileserver".
	// `Run` automatically passes down a `context.Background()` and parses `os.Args[1:]`.
	// If an error is returned, and it is either a `cli.ExitError` or an `*exec.ExitError`, the error's exit code will be used.
	// For any other errors returned, it exits with code 1.
	cli.Run("fileserver", "An HTTP server.",

		// Add an optional positional argument called "root" which will default to ".".
		cli.AddArg("root", "Directory to serve from", cli.SetArgDefault(".")),

		// Add an optional flag called "port" (usage: --port) which will default to 3000.
		cli.AddFlag("port", "Port to run server.", cli.SetFlagDefault(3000)),

		// Register a handler for this command.
		// If no handler is registered, it will simply print help and exit.
		cli.SetHandler(func(ctx context.Context) error {
			// Extract the value of the "root" argument.
			root, _ := cli.ArgValue[string](ctx, "root")

			// Extract the value of the "port" flag.
			port, _ := cli.FlagValue[int](ctx, "port")

			addr := fmt.Sprintf(":%d", port)
			return http.ListenAndServe(addr, http.FileServer(http.Dir(root)))
		}),
	)
}

func Example_kitchen_sink() {
	// Create a new command
	cmd, err := cli.NewCommand("git", "Modern version control.",
		// Set command version
		cli.SetVersion("2.37.0"),

		// Add a "--version" flag with a short flag "-V" for printing the command version
		cli.AddVersionFlag(cli.AddFlagShort('V')),

		// Add a "--help" flag
		cli.AddHelpFlag(

			// Add a short flag "-h" to help
			cli.AddFlagShort('h'),

			// Make this flag inherited by sub-commands
			cli.SetFlagIsInherited(true),
		),

		// Add a hidden "--debug" flag
		cli.AddFlag("debug", "Enable debugging",
			cli.SetFlagDefault(false), // Default parser for flags is cli.StringParser

			// Make it hidden
			cli.SetFlagIsHidden(true),
		),

		// Add a sub-command "clone"
		cli.AddSubCmd("clone", "Clone a repository.",

			// Add a required argument "<url>"
			cli.AddArg("url", "Repository to clone.",

				// Parse it into a *url.URL
				cli.SetArgParser(cli.URLParser),
			),

			// Add optional argument "<dir?>"
			cli.AddArg("dir", "Directory to clone repo into.",

				// Set its default value to "."
				cli.SetArgDefault("."),
			),

			// Add a flag "--verbose"
			cli.AddFlag("verbose", "Be more verbose.",

				// Add a short "-v"
				cli.AddFlagShort('v'),

				// Make it a boolean that defaults to false
				cli.SetFlagDefault(false),
			),
		),
	)
	if err != nil {
		cli.ExitWithError(err)
	}

	// Pass in your `context.Context` and args
	if err := cmd.Run(context.TODO(), os.Args[1:]); err != nil {
		cli.ExitWithError(err)
	}
}
