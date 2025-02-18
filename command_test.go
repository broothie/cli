package cli

import (
	"os"
)

func ExampleNewCommand() {
	command, _ := NewCommand("server", "An http server.",
		AddHelpFlag(AddFlagShort('h')),
		SetVersion("v0.1.0"),
		AddVersionFlag(AddFlagShort('V')),
		AddFlag("port", "Port to run server on.",
			SetFlagDefault(3000),
			AddFlagShort('p'),
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
	//   server [flags] [sub-commands]
	//
	// Sub-commands:
	//   proxy: Proxy requests to another server.
	//
	// Flags:
	//   --help           -h  Print help.                         (type: bool, default: "false")
	//   --version        -V  Print version.                      (type: bool, default: "false")
	//   --port           -p  Port to run server on.              (type: int, default: "3000")
	//   --auth-required      Whether to require authentication.  (type: bool, default: "true")
}
