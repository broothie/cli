package cli

import "github.com/bobg/errors"

func (c *Command) validateInput() error {
	validations := []func() error{
		c.validateArgumentsInput,
	}

	var errs []error
	for _, validation := range validations {
		errs = append(errs, validation())
	}

	return errors.Join(errs...)
}

func (c *Command) validateArgumentsInput() error {
	var errs []error
	for _, argument := range c.arguments {
		errs = append(errs, argument.validateInput())
	}

	return errors.Join(errs...)
}
