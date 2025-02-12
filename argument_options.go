package cli

import "github.com/broothie/option"

func SetArgParser[T any](parser ArgParser[T]) option.Func[*Argument] {
	return func(argument *Argument) (*Argument, error) {
		argument.parser = parser
		return argument, nil
	}
}
