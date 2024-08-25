package cli

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/bobg/errors"
)

type ValueParser func(value string) (any, error)

func StringParser(value string) (any, error) {
	return value, nil
}

func BoolParser(value string) (any, error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing %q as bool", value)
	}

	return b, nil
}

func IntParser(value string) (any, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing %q as int", i)
	}

	return i, nil
}

func Float64Parser(value string) (any, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing %q as float", value)
	}

	return f, nil
}

func URLParser(value string) (any, error) {
	u, err := url.Parse(value)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing %q as url", value)
	}

	return u, nil
}

type Parseable interface {
	string | bool | int | float64 | *url.URL
}

func valueParserFromParseable[T Parseable]() (ValueParser, error) {
	var t T
	switch any(t).(type) {
	case string:
		return StringParser, nil

	case bool:
		return BoolParser, nil

	case int:
		return IntParser, nil

	case float64:
		return Float64Parser, nil

	case *url.URL:
		return URLParser, nil

	default:
		return nil, fmt.Errorf("invalid type %T: not Parseable", t)
	}
}
