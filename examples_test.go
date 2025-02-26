package cli_test

import (
	"context"
	"fmt"
	"net/http"

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

		// Add an optional flag called "port" which will default to 3000.
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
