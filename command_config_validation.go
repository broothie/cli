package cli

import "github.com/bobg/errors"

func (c *Command) validateConfig() error {
	validations := []func() error{
		c.validateNoDuplicateFlags,
		c.validateNoDuplicateArguments,
		c.validateNoDuplicateSubCommands,
		c.validateEitherCommandsOrArguments,
		c.validateVariadicArguments,
	}

	var errs []error
	for _, validation := range validations {
		errs = append(errs, validation())
	}

	return errors.Join(errs...)
}

func (c *Command) validateNoDuplicateFlags() error {
	flags := make(map[string]bool)

	var errs []error
	for _, flag := range c.flags {
		for _, name := range append([]string{flag.name}, flag.aliases...) {
			if flags[name] {
				errs = append(errs, errors.Errorf("duplicate flag %q", name))
			}

			flags[name] = true
		}
	}

	return errors.Join(errs...)
}

func (c *Command) validateNoDuplicateArguments() error {
	arguments := make(map[string]bool)

	var errs []error
	for _, argument := range c.arguments {
		if arguments[argument.name] {
			errs = append(errs, errors.Errorf("duplicate argument %q", argument.name))
		}

		arguments[argument.name] = true
	}

	return errors.Join(errs...)
}

func (c *Command) validateNoDuplicateSubCommands() error {
	commands := make(map[string]bool)

	var errs []error
	for _, command := range c.subCommands {
		if commands[command.name] {
			errs = append(errs, errors.Errorf("duplicate sub-command %q", command.name))
		}

		commands[command.name] = true
	}

	return errors.Join(errs...)
}

func (c *Command) validateEitherCommandsOrArguments() error {
	if len(c.subCommands) > 0 && len(c.arguments) > 0 {
		return errors.Errorf("cannot have both sub-commands and arguments")
	}

	return nil
}

func (c *Command) validateVariadicArguments() error {
	var errs []error

	// Check that only the last argument can be variadic
	for i, argument := range c.arguments {
		if argument.isVariadic() && i != len(c.arguments)-1 {
			errs = append(errs, errors.Errorf("only the last argument can be variadic, but argument %q at position %d is variadic", argument.name, i+1))
		}
	}

	// Check that variadic arguments cannot have default values  
	for _, argument := range c.arguments {
		if argument.isVariadic() && argument.isOptional() {
			errs = append(errs, errors.Errorf("variadic argument %q cannot have a default value", argument.name))
		}
	}

	return errors.Join(errs...)
}
