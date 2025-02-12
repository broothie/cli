package cli

import (
	"context"

	"github.com/bobg/errors"
)

var (
	NotACommandContextError = errors.New("not a command context")
	FlagNotFoundError       = errors.New("flag not found")
	ArgumentNotFoundError   = errors.New("argument not found")
)

func FlagValue(ctx context.Context, name string) (any, error) {
	command, err := commandFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "finding flag %q", name)
	}

	flag, found := command.findFlag(name)
	if !found {
		return nil, errors.Wrapf(FlagNotFoundError, "finding flag %q", name)
	}

	if flag.value != nil {
		return flag.value, nil
	}

	return flag.defaultValue, nil
}

func ArgValue(ctx context.Context, name string) (any, error) {
	command, err := commandFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "finding argument %q", name)
	}

	arg, found := command.findArg(name)
	if !found {
		return nil, errors.Wrapf(ArgumentNotFoundError, "finding argument %q", name)
	}

	return arg.value, nil
}

var commandContextKey struct{}

func (c *Command) onContext(parent context.Context) context.Context {
	return context.WithValue(parent, commandContextKey, c)
}

func commandFromContext(ctx context.Context) (*Command, error) {
	command, ok := ctx.Value(commandContextKey).(*Command)
	if !ok {
		return nil, NotACommandContextError
	}

	return command, nil
}
