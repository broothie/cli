package cli

import (
	"fmt"

	"github.com/samber/lo"
)

type Argument struct {
	name        string
	description string
	valueParser ValueParser

	value any
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
