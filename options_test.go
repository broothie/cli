package cli

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

		assert.NoError(t, err)

		// Command
		assert.Equal(t, "http", httpCommand.name)
		assert.Equal(t, "Run http server", httpCommand.description)
		assert.Equal(t, "v0.1.0", httpCommand.version)
		assert.Nil(t, httpCommand.parent)
		assert.Equal(t, httpCommand, httpCommand.root())
		assert.True(t, httpCommand.isRoot())
		assert.False(t, httpCommand.hasParent())
		assert.Equal(t, "http", httpCommand.qualifiedName())
		assert.Equal(t, "v0.1.0", httpCommand.findVersion())

		// Flags
		assert.NotSliceEmpty(t, httpCommand.flags)

		helpFlag := httpCommand.flags[0]
		assert.Equal(t, "help", helpFlag.name)
		assert.Equal(t, "Print help.", helpFlag.description)
		assert.Equal(t, false, helpFlag.defaultValue)
		assert.Nil(t, helpFlag.value)
		assert.Equal(t, reflect.ValueOf(BoolParser).Pointer(), reflect.ValueOf(helpFlag.parser).Pointer())
		assert.True(t, helpFlag.isBool())
		assert.True(t, helpFlag.isHelp)

		portFlag := httpCommand.flags[1]
		assert.Equal(t, "port", portFlag.name)
		assert.Equal(t, "Port to run server on", portFlag.description)
		assert.DeepEqual(t, []string{"addr"}, portFlag.aliases)
		assert.DeepEqual(t, []rune{'p'}, portFlag.shorts)
		assert.Equal(t, 3000, portFlag.defaultValue)
		assert.Equal(t, "PORT", portFlag.defaultEnvName)
		assert.Nil(t, portFlag.value)
		assert.Equal(t, reflect.ValueOf(IntParser).Pointer(), reflect.ValueOf(portFlag.parser).Pointer())
		assert.False(t, portFlag.isBool())
		assert.False(t, portFlag.isHelp)

		// Sub-command
		assert.NotSliceEmpty(t, httpCommand.subCommands)

		proxySubCommand := httpCommand.subCommands[0]
		assert.Equal(t, "proxy", proxySubCommand.name)
		assert.Equal(t, "Proxy requests", proxySubCommand.description)
		assert.Equal(t, "", proxySubCommand.version)
		assert.DeepEqual(t, []string{"p", "x"}, proxySubCommand.aliases)
		assert.NotNil(t, proxySubCommand.parent)
		assert.Equal(t, httpCommand, proxySubCommand.root())
		assert.False(t, proxySubCommand.isRoot())
		assert.True(t, proxySubCommand.hasParent())
		assert.Equal(t, "http proxy", proxySubCommand.qualifiedName())
		assert.Equal(t, "v0.1.0", proxySubCommand.findVersion())

		// Argument
		assert.NotSliceEmpty(t, proxySubCommand.arguments)

		targetArgument := proxySubCommand.arguments[0]
		assert.Equal(t, "target", targetArgument.name)
		assert.Equal(t, "Target to proxy requests to", targetArgument.description)
		assert.Equal(t, reflect.ValueOf(URLParser).Pointer(), reflect.ValueOf(targetArgument.parser).Pointer())
		assert.Nil(t, targetArgument.value)
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
