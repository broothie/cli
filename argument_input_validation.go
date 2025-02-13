package cli

import "github.com/bobg/errors"

var ArgumentMissingValueError = errors.New("argument missing value")

func (a *Argument) validateInput() error {
	if a.value == nil {
		return errors.Wrapf(ArgumentMissingValueError, "argument %q", a.name)
	}

	return nil
}
