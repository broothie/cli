package cli

import (
	"context"
	"os"

	"github.com/bobg/errors"
)

var HandleExecutionError = errors.New("executing handler")

func FlagValue[T any](ctx context.Context, name string) (T, error) {
	var zero T

	command, err := commandFromContext(ctx)
	if err != nil {
		return zero, err
	}

	flag, found := command.findFlag(name)
	if !found {
		return zero, errors.Wrapf(HandleExecutionError, "finding flag %q", name)
	}

	if flag.value != nil {
		return flag.value.(T), nil
	}

	if flag.defaultEnvName != "" {
		value, err := flag.parser.Parse(os.Getenv(flag.defaultEnvName))
		if err != nil {
			return zero, errors.Wrapf(errors.Join(HandleExecutionError, err), "parsing env for flag %q", name)
		}

		return value.(T), nil
	}

	return flag.defaultValue.(T), nil
}

func ArgValue[T any](ctx context.Context, name string) (T, error) {
	var zero T

	command, err := commandFromContext(ctx)
	if err != nil {
		return zero, errors.Wrapf(err, "finding argument %q", name)
	}

	arg, found := command.findArg(name)
	if !found {
		return zero, errors.Wrapf(HandleExecutionError, "finding argument %q", name)
	}

	if arg.value != nil {
		return arg.value.(T), nil
	}

	return arg.defaultValue.(T), nil
}

func Rest(ctx context.Context) ([]string, error) {
	command, err := commandFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting remaining args")
	}

	return command.rest, nil
}

type commandContextKeyType struct{}

var commandContextKey = commandContextKeyType{}

func (c *Command) onContext(parent context.Context) context.Context {
	return context.WithValue(parent, commandContextKey, c)
}

func commandFromContext(ctx context.Context) (*Command, error) {
	command, ok := ctx.Value(commandContextKey).(*Command)
	if !ok {
		return nil, errors.Wrap(HandleExecutionError, "not a command context")
	}

	return command, nil
}
