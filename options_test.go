package cli

import (
	"context"
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
		assert.NotEqual(t, reflect.ValueOf(helpHandler).Pointer(), reflect.ValueOf(httpCommand.handler).Pointer())
		assert.Nil(t, httpCommand.parent)
		assert.True(t, httpCommand.isRoot())
		assert.False(t, httpCommand.hasParent())
		assert.Equal(t, "http", httpCommand.qualifiedName())
		assert.Equal(t, "v0.1.0", httpCommand.findVersion())

		// Flags
		assert.NotEmpty(t, httpCommand.flags)

		helpFlag := httpCommand.flags[0]
		assert.Equal(t, "help", helpFlag.name)
		assert.Equal(t, "Print help", helpFlag.description)
		assert.Equal(t, false, helpFlag.defaultValue)
		assert.Nil(t, helpFlag.value)
		assert.Equal(t, reflect.ValueOf(BoolParser).Pointer(), reflect.ValueOf(helpFlag.parser).Pointer())
		assert.True(t, helpFlag.isBool())
		assert.True(t, helpFlag.isHelp())

		portFlag := httpCommand.flags[1]
		assert.Equal(t, "port", portFlag.name)
		assert.Equal(t, "Port to run server on", portFlag.description)
		assert.Equal(t, []string{"addr"}, portFlag.aliases)
		assert.Equal(t, []rune{'p'}, portFlag.shorts)
		assert.Equal(t, 3000, portFlag.defaultValue)
		assert.Nil(t, portFlag.value)
		assert.Equal(t, reflect.ValueOf(IntParser).Pointer(), reflect.ValueOf(portFlag.parser).Pointer())
		assert.False(t, portFlag.isBool())
		assert.False(t, portFlag.isHelp())

		// Sub-command
		assert.NotEmpty(t, httpCommand.subCommands)

		proxySubCommand := httpCommand.subCommands[0]
		assert.Equal(t, "proxy", proxySubCommand.name)
		assert.Equal(t, "Proxy requests", proxySubCommand.description)
		assert.Equal(t, "", proxySubCommand.version)
		assert.Equal(t, []string{"p", "x"}, proxySubCommand.aliases)
		assert.Equal(t, reflect.ValueOf(helpHandler).Pointer(), reflect.ValueOf(proxySubCommand.handler).Pointer())
		assert.NotNil(t, proxySubCommand.parent)
		assert.False(t, proxySubCommand.isRoot())
		assert.True(t, proxySubCommand.hasParent())
		assert.Equal(t, "http proxy", proxySubCommand.qualifiedName())
		assert.Equal(t, "v0.1.0", proxySubCommand.findVersion())

		// Argument
		assert.NotEmpty(t, proxySubCommand.arguments)

		targetArgument := proxySubCommand.arguments[0]
		assert.Equal(t, "target", targetArgument.name)
		assert.Equal(t, "Target to proxy requests to", targetArgument.description)
		assert.Equal(t, reflect.ValueOf(URLParser).Pointer(), reflect.ValueOf(targetArgument.parser).Pointer())
		assert.Nil(t, targetArgument.value)
	})
}
