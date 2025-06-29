package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgument_validateConfig(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		arg := &Argument{name: ""}
		err := arg.validateConfig()
		assert.EqualError(t, err, "argument name cannot be empty")
	})

	t.Run("multiple tokens", func(t *testing.T) {
		arg := &Argument{name: "invalid argument name"}
		err := arg.validateConfig()
		assert.EqualError(t, err, `argument name "invalid argument name" must be a single token`)
	})

	t.Run("valid name", func(t *testing.T) {
		arg := &Argument{name: "valid-arg"}
		err := arg.validateConfig()
		require.NoError(t, err)
	})
}
