package cli

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExitError_Error(t *testing.T) {
	err := &ExitError{Code: 2}
	assert.Equal(t, "exit status 2", err.Error())
}

func TestExitCode(t *testing.T) {
	err := ExitCode(3)
	assert.Equal(t, 3, err.Code)
}

func TestExitWithError(t *testing.T) {
	if os.Getenv("TEST_EXIT") == "1" {
		ExitWithError(ExitCode(4))
		return
	}

	// Test ExitError
	cmd := exec.Command(os.Args[0], "-test.run=TestExitWithError")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		assert.Equal(t, 4, exitErr.ExitCode())
	} else {
		t.Errorf("expected ExitError, got %v", err)
	}
}
