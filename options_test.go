package cli

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/broothie/test"
)

func Test_options(t *testing.T) {
	t.Run("kitchen sink", func(t *testing.T) {
		httpCommand, err := NewCommand("http", "Run http server",
			SetVersion("v0.1.0"),
			AddHelpFlag(),
			AddFlag("port", "Port to run server on",
				AddFlagAlias("addr"),
				AddFlagShort('p'),
				SetFlagDefault(3000),
				SetFlagDefaultEnv("PORT"),
			),
			AddSubCmd("proxy", "Proxy requests",
				AddAlias("p"),
				AddAlias("x"),
				AddArg("target", "Target to proxy requests to", SetArgParser(URLParser)),
			),
			SetHandler(func(context.Context) error { return nil }),
		)

		test.Nil(t, err)

		// Command
		test.Equal(t, "http", httpCommand.name)
		test.Equal(t, "Run http server", httpCommand.description)
		test.Equal(t, "v0.1.0", httpCommand.version)
		test.Nil(t, httpCommand.parent)
		test.Equal(t, httpCommand, httpCommand.root())
		test.True(t, httpCommand.isRoot())
		test.False(t, httpCommand.hasParent())
		test.Equal(t, "http", httpCommand.qualifiedName())
		test.Equal(t, "v0.1.0", httpCommand.findVersion())

		// Flags
		test.NotSliceEmpty(t, httpCommand.flags)

		helpFlag := httpCommand.flags[0]
		test.Equal(t, "help", helpFlag.name)
		test.Equal(t, "Print help.", helpFlag.description)
		test.Equal(t, false, helpFlag.defaultValue)
		test.Nil(t, helpFlag.value)
		test.Equal(t, reflect.ValueOf(BoolParser).Pointer(), reflect.ValueOf(helpFlag.parser).Pointer())
		test.True(t, helpFlag.isBool())
		test.True(t, helpFlag.isHelp)

		portFlag := httpCommand.flags[1]
		test.Equal(t, "port", portFlag.name)
		test.Equal(t, "Port to run server on", portFlag.description)
		test.DeepEqual(t, []string{"addr"}, portFlag.aliases)
		test.DeepEqual(t, []rune{'p'}, portFlag.shorts)
		test.Equal(t, 3000, portFlag.defaultValue)
		test.Equal(t, "PORT", portFlag.defaultEnvName)
		test.Nil(t, portFlag.value)
		test.Equal(t, reflect.ValueOf(IntParser).Pointer(), reflect.ValueOf(portFlag.parser).Pointer())
		test.False(t, portFlag.isBool())
		test.False(t, portFlag.isHelp)

		// Sub-command
		test.NotSliceEmpty(t, httpCommand.subCommands)

		proxySubCommand := httpCommand.subCommands[0]
		test.Equal(t, "proxy", proxySubCommand.name)
		test.Equal(t, "Proxy requests", proxySubCommand.description)
		test.Equal(t, "", proxySubCommand.version)
		test.DeepEqual(t, []string{"p", "x"}, proxySubCommand.aliases)
		test.NotNil(t, proxySubCommand.parent)
		test.Equal(t, httpCommand, proxySubCommand.root())
		test.False(t, proxySubCommand.isRoot())
		test.True(t, proxySubCommand.hasParent())
		test.Equal(t, "http proxy", proxySubCommand.qualifiedName())
		test.Equal(t, "v0.1.0", proxySubCommand.findVersion())

		// Argument
		test.NotSliceEmpty(t, proxySubCommand.arguments)

		targetArgument := proxySubCommand.arguments[0]
		test.Equal(t, "target", targetArgument.name)
		test.Equal(t, "Target to proxy requests to", targetArgument.description)
		test.Equal(t, reflect.ValueOf(URLParser).Pointer(), reflect.ValueOf(targetArgument.parser).Pointer())
		test.Nil(t, targetArgument.value)
	})
}

func ExampleSetVersion() {
	command, _ := NewCommand("server", "An http server.",
		SetVersion("v0.1.0"),
	)

	command.renderHelp(os.Stdout)
	// Output:
	// server v0.1.0: An http server.
	//
	// Usage:
	//   server
}

func ExampleSetHandler() {
	command, _ := NewCommand("server", "An http server.",
		SetHandler(func(ctx context.Context) error {
			fmt.Println("running server")
			return nil
		}),
	)

	command.Run(context.Background(), nil)
	// Output:
	// running server
}

func ExampleAddSubCmd() {
	command, _ := NewCommand("server", "An http server.",
		AddSubCmd("start", "Start the server"),
	)

	command.renderHelp(os.Stdout)
	// Output:
	// server: An http server.
	//
	// Usage:
	//   server [sub-command]
	//
	// Sub-command:
	//   start: Start the server
}

func ExampleAddFlag() {
	command, _ := NewCommand("server", "An http server.",
		AddFlag("port", "Port to run server on",
			AddFlagShort('p'),
			SetFlagDefault(3000),
		),
	)

	command.renderHelp(os.Stdout)
	// Output:
	// server: An http server.
	//
	// Usage:
	//   server [flags]
	//
	// Flags:
	//   --port  -p  Port to run server on  (type: int, default: "3000")
}

func ExampleAddFlag_with_env() {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	command, _ := NewCommand("server", "An http server.",
		AddFlag("port", "Port to run server on",
			AddFlagShort('p'),
			SetFlagDefault(3000),
			SetFlagDefaultEnv("PORT"),
		),
	)

	command.renderHelp(os.Stdout)
	// Output:
	// server: An http server.
	//
	// Usage:
	//   server [flags]
	//
	// Flags:
	//   --port  -p  Port to run server on  (type: int, default: $PORT, "3000")
}

func ExampleAddArg() {
	command, _ := NewCommand("server", "An http server.",
		AddArg("port", "Port to run server on", SetArgParser(IntParser)),
	)

	command.renderHelp(os.Stdout)
	// Output:
	// server: An http server.
	//
	// Usage:
	//   server <port>
	//
	// Arguments:
	//   <port>  Port to run server on  (type: int)
}

func ExampleAddHelpFlag() {
	command, _ := NewCommand("server", "An http server.",
		AddHelpFlag(AddFlagShort('h')),
	)

	command.renderHelp(os.Stdout)
	// Output:
	// server: An http server.
	//
	// Usage:
	//   server [flags]
	//
	// Flags:
	//   --help  -h  Print help.  (type: bool, default: "false")
}

func ExampleAddVersionFlag() {
	command, _ := NewCommand("server", "An http server.",
		SetVersion("v0.1.0"),
		AddVersionFlag(AddFlagShort('V')),
	)

	command.renderHelp(os.Stdout)
	// Output:
	// server v0.1.0: An http server.
	//
	// Usage:
	//   server [flags]
	//
	// Flags:
	//   --version  -V  Print version.  (type: bool, default: "false")
}
