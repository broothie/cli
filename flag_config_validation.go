package cli

import (
	"strings"

	"github.com/bobg/errors"
)

func (f *Flag) validateConfig() error {
	if len(strings.Fields(f.name)) > 1 {
		return errors.Errorf("flag name %q must be a single token", f.name)
	}

	if f.name == "" {
		return errors.New("flag name cannot be empty")
	}

	return nil
}
