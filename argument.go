package cli

import (
	"fmt"

	"github.com/bobg/errors"
	"github.com/broothie/option"
	"github.com/samber/lo"
)

type Argument struct {
	name         string
	description  string
	parser       argParser
	defaultValue any

	value any
}

func newArgument(name, description string, options ...option.Option[*Argument]) (*Argument, error) {
	baseArgument := &Argument{
		name:        name,
		description: description,
		parser:      NewArgParser(StringParser),
	}

	argument, err := option.Apply(baseArgument, options...)
	if err != nil {
		return nil, errors.Wrapf(err, "building argument %q", name)
	}

	if err := argument.validateConfig(); err != nil {
		return nil, errors.Wrapf(err, "invalid argument %q", name)
	}

	return argument, nil
}

func (a *Argument) isRequired() bool {
	return a.defaultValue == nil
}

func (a *Argument) isOptional() bool {
	return !a.isRequired()
}

func (a *Argument) inBrackets() string {
	if a.isOptional() {
		return fmt.Sprintf("[<%s>]", a.name)
	}

	return fmt.Sprintf("<%s>", a.name)
}

func (c *Command) findArg(name string) (*Argument, bool) {
	if arg, found := lo.Find(c.arguments, func(argument *Argument) bool { return argument.name == name }); found {
		return arg, true
	}

	return nil, false
}
