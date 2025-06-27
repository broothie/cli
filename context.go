package cli

import (
	"context"
	"os"

	"github.com/bobg/errors"
)

var (
	NotACommandContextError = errors.New("not a command context")
	FlagNotFoundError       = errors.New("flag not found")
	ArgumentNotFoundError   = errors.New("argument not found")
)

func FlagValue[T any](ctx context.Context, name string) (T, error) {
	var zero T

	command, err := commandFromContext(ctx)
	if err != nil {
		return zero, errors.Wrapf(err, "finding flag %q", name)
	}

	flag, found := command.findFlag(name)
	if !found {
		return zero, errors.Wrapf(FlagNotFoundError, "finding flag %q", name)
	}

	if flag.value != nil {
		return flag.value.(T), nil
	}

	if flag.defaultEnvName != "" {
		value, err := flag.parser.Parse(os.Getenv(flag.defaultEnvName))
		if err != nil {
			return zero, err
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
		return zero, errors.Wrapf(ArgumentNotFoundError, "finding argument %q", name)
	}

	if arg.value != nil {
		return arg.value.(T), nil
	}

	return arg.defaultValue.(T), nil
}

func VariadicArgValue[T any](ctx context.Context, name string) ([]T, error) {
	command, err := commandFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "finding variadic argument %q", name)
	}

	arg, found := command.findArg(name)
	if !found {
		return nil, errors.Wrapf(ArgumentNotFoundError, "finding variadic argument %q", name)
	}

	if !arg.isVariadic() {
		return nil, errors.Errorf("argument %q is not variadic", name)
	}

	if arg.value == nil {
		return []T{}, nil
	}

	// Convert []any to []T
	anySlice := arg.value.([]any)
	result := make([]T, len(anySlice))
	for i, v := range anySlice {
		if val, ok := v.(T); ok {
			result[i] = val
		} else {
			return nil, errors.Errorf("invalid type in variadic argument %q at index %d: expected %T, got %T", name, i, *new(T), v)
		}
	}

	return result, nil
}

type commandContextKeyType struct{}

var commandContextKey = commandContextKeyType{}

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
