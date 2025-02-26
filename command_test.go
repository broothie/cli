package cli

import (
	"context"
	"fmt"
	"os"
)

func ExampleNewCommand() {
	command, _ := NewCommand("server", "An http server.",
		SetVersion("v0.1.0"),
		AddVersionFlag(AddFlagShort('V')),
		AddHelpFlag(AddFlagShort('h')),
		AddFlag("port", "Port to run server on.",
			AddFlagShort('p'),
			SetFlagDefault(3000),
			SetFlagDefaultEnv("PORT"),
		),
		AddFlag("auth-required", "Whether to require authentication.",
			SetFlagDefault(true),
		),
		AddSubCmd("proxy", "Proxy requests to another server.",
			AddArg("target", "Target server to proxy requests to.",
				SetArgParser(URLParser),
			),
		),
	)

	command.renderHelp(os.Stdout)
	// Output:
	// server v0.1.0: An http server.
	//
	// Usage:
	//   server [flags] [sub-command]
	//
	// Sub-command:
	//   proxy: Proxy requests to another server.
	//
	// Flags:
	//   --version        -V  Print version.                      (type: bool, default: "false")
	//   --help           -h  Print help.                         (type: bool, default: "false")
	//   --port           -p  Port to run server on.              (type: int, default: $PORT, "3000")
	//   --auth-required      Whether to require authentication.  (type: bool, default: "true")
}

func ExampleRun() {
	Run("echo", "Echo the arguments.",
		AddArg("arg", "The argument to echo.", SetArgDefault("hello")),
		SetHandler(func(ctx context.Context) error {
			arg, _ := ArgValue[string](ctx, "arg")

			fmt.Println(arg)
			return nil
		}),
	)

	// Output:
	// hello
}
