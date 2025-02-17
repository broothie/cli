package cli

import "github.com/broothie/option"

// SetArgParser sets the parser of the argument.
func SetArgParser[T any](parser ArgParser[T]) option.Func[*Argument] {
	return func(argument *Argument) (*Argument, error) {
		argument.parser = parser
		return argument, nil
	}
}
