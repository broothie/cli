package cli

import (
	"net/url"
	"strconv"
	"time"

	"github.com/bobg/errors"
)

var NotParseableError = errors.New("type not parseable")

// Parseable is a type that can be parsed from a string.
type Parseable interface {
	string | bool | int | float64 | time.Time | time.Duration | *url.URL
}

// StringParser parses a string into a string.
func StringParser(s string) (string, error) {
	return s, nil
}

// BoolParser parses a string into a bool.
func BoolParser(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// IntParser parses a string into an int.
func IntParser(s string) (int, error) {
	return strconv.Atoi(s)
}

// Float64Parser parses a string into a float64.
func Float64Parser(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// TimeParser parses a string into a time.Time.
func TimeParser(s string) (time.Time, error) {
	return TimeLayoutParser(time.RFC3339)(s)
}

// TimeLayoutParser parses a string into a time.Time using a specific time layout.
func TimeLayoutParser(timeLayout string) ArgParser[time.Time] {
	return func(s string) (time.Time, error) {
		return time.Parse(timeLayout, s)
	}
}

// DurationParser parses a string into a time.Duration.
func DurationParser(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// URLParser parses a string into a *url.URL.
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
