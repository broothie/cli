package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/samber/lo"
)

func TestCommand_flagsUpToRoot(t *testing.T) {
	command, err := NewCommand("test", "test command",
		AddFlag("top-uninherited", "top uninherited"),
		AddFlag("top-inherited", "top inherited", SetFlagIsInherited(true)),
		AddSubCmd("sub-command", "sub-command",
			AddFlag("flag", "flag"),
		),
	)

	require.NoError(t, err)

	flags := command.subCommands[0].flagsUpToRoot()
	flagNames := lo.Map(flags, func(flag *Flag, _ int) string { return flag.name })
	assert.Contains(t, flagNames, "top-inherited")
	assert.Contains(t, flagNames, "flag")
	assert.NotContains(t, flagNames, "top-uninherited")
}
