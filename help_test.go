package cli

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/broothie/test"
)

type CustomType struct {
	Field string
}

func (c CustomType) String() string {
	return fmt.Sprintf("field=%q", c.Field)
}

func TestCommand_renderHelp(t *testing.T) {
	t.Run("with args, no command", func(t *testing.T) {
		command, err := NewCommand("test", "test command",
			SetVersion("v1.2.3-rc10"),
			AddHelpFlag(),
			AddFlag("some-flag", "some flag",
				AddFlagAlias("a-flag"),
				AddFlagShort('s'),
				SetFlagDefaultAndParser(CustomType{Field: "field default"}, func(s string) (CustomType, error) { return CustomType{Field: s}, nil }),
			),
			AddFlag("hidden-flag", "some hidden flag", SetFlagIsHidden(true)),
			AddArg("some-arg", "some arg",
				SetArgParser(TimeLayoutParser(time.RubyDate)),
			),
		)

		test.Nil(t, err)

		buffer := new(bytes.Buffer)
		test.Nil(t, command.renderHelp(buffer))

		test.Equal(t,
			heredoc.Doc(`
				test v1.2.3-rc10: test command
				
				Usage:
				  test [flags] <some-arg>
				
				Arguments:
				  <some-arg>  some arg  (type: time.Time)
				
				Flags:
				  --help                    Print help.  (type: bool, default: "false")
				  --some-flag --a-flag  -s  some flag    (type: cli.CustomType, default: "field=\"field default\"")
				
			`),
			buffer.String(),
		)
	})

	t.Run("with sub-command, no args", func(t *testing.T) {
		command, err := NewCommand("test", "test command",
			SetVersion("v1.2.3-rc10"),
			AddHelpFlag(),
			AddFlag("some-flag", "some flag",
				AddFlagAlias("a-flag"),
				AddFlagShort('s'),
				SetFlagDefaultAndParser(CustomType{Field: "field default"}, func(s string) (CustomType, error) { return CustomType{Field: s}, nil }),
			),
			AddSubCmd("some-command", "some command"),
		)

		test.Nil(t, err)

		buffer := new(bytes.Buffer)
		test.Nil(t, command.renderHelp(buffer))

		test.Equal(t,
			heredoc.Doc(`
				test v1.2.3-rc10: test command
				
				Usage:
				  test [flags] [sub-commands]
				
				Sub-commands:
				  some-command: some command
				
				Flags:
				  --help                    Print help.  (type: bool, default: "false")
				  --some-flag --a-flag  -s  some flag    (type: cli.CustomType, default: "field=\"field default\"")
				
			`),
			buffer.String(),
		)
	})

	t.Run("sub-command", func(t *testing.T) {
		command, err := NewCommand("test", "test command",
			SetVersion("v1.2.3-rc10"),
			AddHelpFlag(),
			AddFlag("inherited", "inherited", SetFlagIsInherited(true)),
			AddSubCmd("some-command", "some command",
				AddFlag("some-flag", "some flag"),
				AddArg("some-arg", "some arg"),
			),
		)

		test.Nil(t, err)

		buffer := new(bytes.Buffer)
		test.Nil(t, command.subCommands[0].renderHelp(buffer))

		test.Equal(t,
			heredoc.Doc(`
				test v1.2.3-rc10: test command
				
				Usage:
				  test some-command [flags] <some-arg>
				
				Arguments:
				  <some-arg>  some arg  (type: string)
				
				Flags:
				  --some-flag    some flag  (type: string, default: "")
				  --inherited    inherited  (type: string, default: "")
				
			`),
			buffer.String(),
		)
	})
}
