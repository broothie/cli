package cli

import "fmt"

type ExitError struct {
	Code int
}

func (e ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

func ExitCode(code int) ExitError {
	return ExitError{Code: code}
}
