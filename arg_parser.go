package cli

type argParser interface {
	Type() any
	Parse(string) (any, error)
}

// ArgParser is a function that parses a string into a value of type T.
type ArgParser[T any] func(string) (T, error)

// NewArgParser creates a new ArgParser.
func NewArgParser[T any](f ArgParser[T]) ArgParser[T] {
	return f
}

func (ArgParser[T]) Type() any {
	var t T
	return t
}

func (p ArgParser[T]) Parse(s string) (any, error) {
	return p(s)
}
