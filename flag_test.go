package cli

import (
	"testing"

	"github.com/broothie/test"
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

	test.NoError(t, err)

	flags := command.subCommands[0].flagsUpToRoot()
	flagNames := lo.Map(flags, func(flag *Flag, _ int) string { return flag.name })
	test.Contains(t, flagNames, "top-inherited")
	test.Contains(t, flagNames, "flag")
	test.NotContains(t, flagNames, "top-uninherited")
}
