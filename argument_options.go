package cli

import "github.com/broothie/option"

// SetArgParser sets the parser of the argument.
func SetArgParser[T any](parser ArgParser[T]) option.Func[*Argument] {
	return func(argument *Argument) (*Argument, error) {
		argument.parser = parser
		return argument, nil
	}
}

// SetArgDefault sets the default value of the argument.
func SetArgDefault[T Parseable](defaultValue T) option.Func[*Argument] {
	return func(argument *Argument) (*Argument, error) {
		argParser, err := argParserFromParseable[T]()
		if err != nil {
			return nil, err
		}

		argument.parser = argParser
		argument.defaultValue = defaultValue
		return argument, nil
	}
}
