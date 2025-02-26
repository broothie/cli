package cli

import (
	"testing"

	"github.com/broothie/test"
)

func TestArgument_validateInput(t *testing.T) {
	arg, err := newArgument("test-arg", "Test arg.")
	test.MustNoError(t, err)

	test.ErrorMessageIs(t, arg.validateInput(), `argument "test-arg": argument missing value`)

	arg.value = "something"
	test.NoError(t, arg.validateInput())
}
