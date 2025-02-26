package cli

import (
	"fmt"

	"github.com/bobg/errors"
	"github.com/broothie/option"
	"github.com/samber/lo"
)

const helpFlagName = "help"

type Flag struct {
	name           string
	description    string
	aliases        []string
	shorts         []rune
	isHelp         bool
	isVersion      bool
	isHidden       bool
	isInherited    bool
	parser         argParser
	defaultEnvName string
	defaultValue   any

	value any
}

func newFlag(name, description string, options ...option.Option[*Flag]) (*Flag, error) {
	baseFlag := &Flag{
		name:         name,
		description:  description,
		parser:       NewArgParser(StringParser),
		defaultValue: "",
	}

	flag, err := option.Apply(baseFlag, options...)
	if err != nil {
		return nil, errors.Wrapf(err, "building flag %q", name)
	}

	return flag, nil
}

func (f *Flag) isBool() bool {
	return isBoolParser(f.parser)
}

func (c *Command) findFlag(name string) (*Flag, bool) {
	return c.findFlagUpToRoot(func(flag *Flag) bool { return flag.name == name })
}

func (c *Command) findFlagUpToRoot(predicate func(*Flag) bool) (*Flag, bool) {
	for current := c; current != nil; current = current.parent {
		currentIsSelf := current == c

		flags := current.flags
		if !currentIsSelf {
			flags = lo.Filter(flags, func(flag *Flag, _ int) bool { return flag.isInherited })
		}

		flag, found := lo.Find(flags, predicate)
		if found {
			return flag, true
		}
	}

	return nil, false
}

func (c *Command) flagsUpToRoot() []*Flag {
	flags := c.flags
	flagSet := lo.Associate(flags, func(flag *Flag) (string, bool) { return flag.name, true })

	for current := c.parent; current != nil; current = current.parent {
		flags = append(flags, lo.Filter(current.flags, func(flag *Flag, _ int) bool {
			defer func() { flagSet[flag.name] = true }()

			return flag.isInherited && !flagSet[flag.name]
		})...)
	}

	return flags
}

func dashifyShort(short rune) string {
	return fmt.Sprintf("-%c", short)
}

func isBoolParser(parser argParser) bool {
	_, isBool := parser.Type().(bool)
	return isBool
}
