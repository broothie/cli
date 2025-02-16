package cli

import (
	"github.com/bobg/errors"
	"github.com/broothie/option"
)

func AddFlagAlias(alias string) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.aliases = append(flag.aliases, alias)
		return flag, nil
	}
}

func AddFlagShort(short rune) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.shorts = append(flag.shorts, short)
		return flag, nil
	}
}

func SetFlagHidden() option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.hidden = true
		return flag, nil
	}
}

func SetFlagDefault[T Parseable](defaultValue T) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		argParser, err := argParserFromParseable[T]()
		if err != nil {
			return nil, errors.Wrapf(err, "")
		}

		flag.parser = argParser
		flag.defaultValue = defaultValue
		return flag, nil
	}
}

func SetFlagDefaultAndParser[T any](defaultValue T, argParser ArgParser[T]) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.parser = argParser
		flag.defaultValue = defaultValue
		return flag, nil
	}
}
