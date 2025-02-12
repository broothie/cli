package cli

import (
	"fmt"

	"github.com/bobg/errors"
	"github.com/broothie/option"
	"github.com/samber/lo"
)

type Argument struct {
	name        string
	description string
	parser      argParser

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

	return argument, nil
}

func (a *Argument) inBrackets() string {
	return fmt.Sprintf("<%s>", a.name)
}

func (c *Command) findArg(name string) (*Argument, bool) {
	if arg, found := lo.Find(c.arguments, func(argument *Argument) bool { return argument.name == name }); found {
		return arg, true
	}

	if c.hasParent() {
		return c.parent.findArg(name)
	}

	return nil, false
}
