package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand_root(t *testing.T) {
	command, err := NewCommand("root", "", AddSubCmd("child", ""))
	require.NoError(t, err)

	require.NotEmpty(t, command.subCommands)
	subCommand := command.subCommands[0]

	assert.True(t, command.isRoot())
	assert.False(t, subCommand.isRoot())
}

func TestCommand_hasParent(t *testing.T) {
	command, err := NewCommand("root", "", AddSubCmd("child", ""))
	require.NoError(t, err)

	require.NotEmpty(t, command.subCommands)
	subCommand := command.subCommands[0]

	assert.False(t, command.hasParent())
	assert.True(t, subCommand.hasParent())
}

func TestCommand_qualifiedName(t *testing.T) {
	command, err := NewCommand("root", "",
		AddSubCmd("child", "",
			AddSubCmd("grandchild", ""),
		),
	)
	require.NoError(t, err)

	require.NotEmpty(t, command.subCommands)
	childCommand := command.subCommands[0]

	require.NotEmpty(t, childCommand.subCommands)
	grandChildCommand := childCommand.subCommands[0]

	assert.Equal(t, command.qualifiedName(), "root")
	assert.Equal(t, childCommand.qualifiedName(), "root child")
	assert.Equal(t, grandChildCommand.qualifiedName(), "root child grandchild")
}

func TestCommand_findVersion(t *testing.T) {
	command, err := NewCommand("root", "",
		SetVersion("v0.1.0"),
		AddSubCmd("child", "",
			AddSubCmd("grandchild", "",
				SetVersion("v0.3.0"),
			),
		),
	)

	assert.NoError(t, err)

	require.NotEmpty(t, command.subCommands)
	childCommand := command.subCommands[0]

	require.NotEmpty(t, childCommand.subCommands)
	grandChildCommand := childCommand.subCommands[0]

	assert.Equal(t, command.findVersion(), "v0.1.0")
	assert.Equal(t, childCommand.findVersion(), "v0.1.0")
	assert.Equal(t, grandChildCommand.findVersion(), "v0.3.0")
}
