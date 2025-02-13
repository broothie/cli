package cli

import (
	"github.com/bobg/errors"
)

func (c *Command) validateConfig() error {
	validations := []func() error{
		c.validateNoDuplicateSubCommands,
		c.validateEitherCommandsOrArguments,
	}

	var errs []error
	for _, validation := range validations {
		errs = append(errs, validation())
	}

	return errors.Join(errs...)
}

func (c *Command) validateNoDuplicateSubCommands() error {
	set := make(map[string]bool)

	var errs []error
	for _, command := range c.subCommands {
		if set[command.name] {
			errs = append(errs, errors.Errorf("duplicate sub-command %q", command.name))
		} else {
			set[command.name] = true
		}
	}

	return errors.Join(errs...)
}

func (c *Command) validateEitherCommandsOrArguments() error {
	if len(c.subCommands) > 0 && len(c.arguments) > 0 {
		return errors.Errorf("cannot have both sub-commands and arguments")
	}

	return nil
}
