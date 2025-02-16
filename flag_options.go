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

func SetFlagIsHidden(isHidden bool) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.isHidden = isHidden
		return flag, nil
	}
}

func SetFlagIsInherited(isInherited bool) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.isInherited = isInherited
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
