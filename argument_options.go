package cli

import "github.com/broothie/option"

func SetArgParser(valueParser ValueParser) option.Func[*Argument] {
	return func(argument *Argument) (*Argument, error) {
		argument.valueParser = valueParser
		return argument, nil
	}
}
