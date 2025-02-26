package cli

import (
	"testing"

	"github.com/broothie/test"
)

func TestArgument_validateConfig(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		arg := &Argument{name: ""}
		err := arg.validateConfig()
		test.ErrorMessageIs(t, err, "argument name cannot be empty")
	})

	t.Run("multiple tokens", func(t *testing.T) {
		arg := &Argument{name: "invalid argument name"}
		err := arg.validateConfig()
		test.ErrorMessageIs(t, err, `argument name "invalid argument name" must be a single token`)
	})

	t.Run("valid name", func(t *testing.T) {
		arg := &Argument{name: "valid-arg"}
		err := arg.validateConfig()
		test.NoError(t, err)
	})
}
