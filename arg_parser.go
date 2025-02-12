package cli

import (
	"net/url"
	"strconv"
)

type argParser interface {
	Type() any
	Parse(string) (any, error)
}

type ArgParser[T any] func(string) (T, error)

func NewArgParser[T any](f ArgParser[T]) ArgParser[T] {
	return f
}

func (ArgParser[T]) Type() any {
	var t T
	return t
}

func (p ArgParser[_]) Parse(s string) (any, error) {
	return p(s)
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

func URLParser(s string) (*url.URL, error) {
	return url.Parse(s)
}
