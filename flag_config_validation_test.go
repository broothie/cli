package cli

import (
	"testing"

	"github.com/broothie/test"
)

func TestFlag_validateConfig(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		flag := &Flag{name: ""}
		err := flag.validateConfig()
		test.ErrorMessageIs(t, err, "flag name cannot be empty")
	})

	t.Run("multiple tokens", func(t *testing.T) {
		flag := &Flag{name: "invalid flag name"}
		err := flag.validateConfig()
		test.ErrorMessageIs(t, err, `flag name "invalid flag name" must be a single token`)
	})

	t.Run("valid name", func(t *testing.T) {
		flag := &Flag{name: "valid-flag"}
		err := flag.validateConfig()
		test.NoError(t, err)
	})
}
