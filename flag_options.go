package cli

import (
	"github.com/bobg/errors"
	"github.com/broothie/option"
)

// AddFlagAlias adds an alias to the flag.
func AddFlagAlias(alias string) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.aliases = append(flag.aliases, alias)
		return flag, nil
	}
}

// AddFlagShort adds a short flag to the flag.
func AddFlagShort(short rune) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.shorts = append(flag.shorts, short)
		return flag, nil
	}
}

// SetFlagIsHidden controls whether the flag is hidden from the help message.
func SetFlagIsHidden(isHidden bool) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.isHidden = isHidden
		return flag, nil
	}
}

// SetFlagIsInherited controls whether the flag is inherited by child commands.
func SetFlagIsInherited(isInherited bool) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.isInherited = isInherited
		return flag, nil
	}
}

// SetFlagDefault sets the default value of the flag.
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

// SetFlagDefaultAndParser sets the default value and parser of the flag.
func SetFlagDefaultAndParser[T any](defaultValue T, argParser ArgParser[T]) option.Func[*Flag] {
	return func(flag *Flag) (*Flag, error) {
		flag.parser = argParser
		flag.defaultValue = defaultValue
		return flag, nil
	}
}
