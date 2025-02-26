package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bobg/errors"
)

// ExitError is an error that causes the program to exit with a given status code.
type ExitError struct {
	Code int
}

// Error implements the error interface.
func (e ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

// ExitCode returns an ExitError with the given code.
func ExitCode(code int) *ExitError {
	return &ExitError{Code: code}
}

// ExitWithError exits the program with an error.
func ExitWithError(err error) {
	fmt.Println(err)

	if exitErr := new(ExitError); errors.As(err, &exitErr) {
		os.Exit(exitErr.Code)
	} else if exitErr := new(exec.ExitError); errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode())
	} else if err != nil {
		os.Exit(1)
	}
}
