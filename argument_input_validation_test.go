package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgument_validateInput(t *testing.T) {
	arg, err := newArgument("test-arg", "Test arg.")
	assert.MustNoError(t, err)

	assert.ErrorMessageIs(t, arg.validateInput(), `argument "test-arg": argument missing value`)

	arg.value = "something"
	assert.NoError(t, arg.validateInput())
}
