package cli

import (
	"testing"

	"github.com/broothie/option"
	"github.com/stretchr/testify/assert"
)

func TestCommand_config_validations(t *testing.T) {
	type TestCase struct {
		commandOptions option.Options[*Command]
		expectedError  string
	}

	testCases := map[string]TestCase{
		"validateNoDuplicateFlags": {
			commandOptions: option.NewOptions(
				AddFlag("some-flag", "some flag"),
				AddFlag("some-flag", "some flag"),
			),
			expectedError: `invalid command "test": duplicate flag "some-flag"`,
		},
		"validateNoDuplicateArguments": {
			commandOptions: option.NewOptions(
				AddArg("some-arg", "some arg"),
				AddArg("some-arg", "some arg"),
			),
			expectedError: `invalid command "test": duplicate argument "some-arg"`,
		},
		"validateNoDuplicateSubCommands": {
			commandOptions: option.NewOptions(
				AddSubCmd("some-sub-command", "some sub-command"),
				AddSubCmd("some-sub-command", "some sub-command"),
			),
			expectedError: `invalid command "test": duplicate sub-command "some-sub-command"`,
		},
		"validateEitherCommandsOrArguments": {
			commandOptions: option.NewOptions(
				AddArg("some-arg", "some arg"),
				AddSubCmd("some-sub-command", "some sub-command"),
			),
			expectedError: `invalid command "test": cannot have both sub-commands and arguments`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := NewCommand("test", "test command", testCase.commandOptions...)
			assert.ErrorMessageIs(t, err, testCase.expectedError)
		})
	}
}
