package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bobg/errors"
)

type ExitError struct {
	Code int
}

func (e ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

func ExitCode(code int) ExitError {
	return ExitError{Code: code}
}

func ExitWithError(err error) {
	fmt.Println(err)

	if exitErr := new(ExitError); errors.As(err, &exitErr) {
		os.Exit(exitErr.Code)
	} else if exitErr := new(exec.ExitError); errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode())
	} else {
		os.Exit(1)
	}
}
