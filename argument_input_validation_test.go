package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgument_validateInput(t *testing.T) {
	arg, err := newArgument("test-arg", "Test arg.")
	require.NoError(t, err)

	assert.EqualError(t, arg.validateInput(), `argument "test-arg": argument missing value`)

	arg.value = "something"
	assert.NoError(t, arg.validateInput())
}
