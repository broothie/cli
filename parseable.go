package cli

import (
	"net/url"
	"strconv"
	"time"

	"github.com/bobg/errors"
)

var NotParseableError = errors.New("type not parseable")

type Parseable interface {
	string | bool | int | float64 | time.Time | time.Duration | *url.URL
}

func StringParser(s string) (string, error) {
	return s, nil
}

func BoolParser(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func IntParser(s string) (int, error) {
	return strconv.Atoi(s)
}

func Float64Parser(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func TimeParser(s string) (time.Time, error) {
	return TimeLayoutParser(time.RFC3339)(s)
}

func TimeLayoutParser(timeLayout string) ArgParser[time.Time] {
	return func(s string) (time.Time, error) {
		return time.Parse(timeLayout, s)
	}
}

func DurationParser(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

func URLParser(s string) (*url.URL, error) {
	return url.Parse(s)
}

func argParserFromParseable[T Parseable]() (argParser, error) {
	var t T
	switch any(t).(type) {
	case string:
		return NewArgParser(StringParser), nil

	case bool:
		return NewArgParser(BoolParser), nil

	case int:
		return NewArgParser(IntParser), nil

	case float64:
		return NewArgParser(Float64Parser), nil

	case time.Time:
		return NewArgParser(TimeParser), nil

	case time.Duration:
		return NewArgParser(DurationParser), nil

	case *url.URL:
		return NewArgParser(URLParser), nil

	default:
		return nil, errors.Wrapf(NotParseableError, "type %T", t)
	}
}
