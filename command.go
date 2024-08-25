package cli

import (
	"context"
	"os"
	"strings"

	"github.com/bobg/errors"
	"github.com/broothie/option"
)

type Handler func(ctx context.Context) error

type Command struct {
	name        string
	description string
	version     string
	aliases     []string
	parent      *Command
	subCommands []*Command
	flags       []*Flag
	arguments   []*Argument
	handler     Handler
}

func New(name, description string, options ...option.Option[*Command]) (*Command, error) {
	command, err := option.Apply(&Command{name: name, description: description, handler: helpHandler}, options...)
	if err != nil {
		return nil, errors.Wrapf(err, "building command %q", name)
	}

	return command, nil
}

func Run(name, description string, options ...option.Option[*Command]) error {
	command, err := New(name, description, options...)
	if err != nil {
		return err
	}

	return command.Run(context.Background(), os.Args[1:])
}

func (c *Command) Run(ctx context.Context, rawArgs []string) error {
	return c.newParser(rawArgs).parse(ctx)
}

func (c *Command) runHandler(ctx context.Context) error {
	return c.handler(c.onContext(ctx))
}

func (c *Command) root() *Command {
	if c.isRoot() {
		return c
	}

	return c.parent.root()
}

func (c *Command) isRoot() bool {
	return c.parent == nil
}

func (c *Command) hasParent() bool {
	return !c.isRoot()
}

func (c *Command) qualifiedName() string {
	if c.hasParent() {
		return strings.Join([]string{c.parent.qualifiedName(), c.name}, " ")
	}

	return c.name
}

func (c *Command) findVersion() string {
	if c.version == "" && c.hasParent() {
		return c.parent.findVersion()
	}

	return c.version
}

var commandContextKey struct{}

func (c *Command) onContext(parent context.Context) context.Context {
	return context.WithValue(parent, commandContextKey, c)
}

func commandFromContext(ctx context.Context) *Command {
	return ctx.Value(commandContextKey).(*Command)
}
