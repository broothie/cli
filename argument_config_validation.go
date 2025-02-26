package cli

import (
	"strings"

	"github.com/bobg/errors"
)

func (a *Argument) validateConfig() error {
	if len(strings.Fields(a.name)) > 1 {
		return errors.Errorf("argument name %q must be a single token", a.name)
	}

	if a.name == "" {
		return errors.New("argument name cannot be empty")
	}

	return nil
}
