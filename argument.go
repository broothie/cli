package cli

import (
	"context"
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

func ArgValue(ctx context.Context, name string) any {
	arg, found := commandFromContext(ctx).findArg(name)
	if !found {
		panic(fmt.Sprintf("no arg with name %q", name))
	}

	return arg.value
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
