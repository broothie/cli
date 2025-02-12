package cli

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
