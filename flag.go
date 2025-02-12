package cli

import (
	"fmt"

	"github.com/bobg/errors"
	"github.com/broothie/option"
	"github.com/samber/lo"
)

const helpFlagName = "help"

type Flag struct {
	name         string
	description  string
	aliases      []string
	shorts       []rune
	parser       argParser
	defaultValue any

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

func (f *Flag) isHelp() bool {
	return f.name == helpFlagName && f.isBool()
}

func (f *Flag) isBool() bool {
	return isBoolParser(f.parser)
}

func (c *Command) findFlag(name string) (*Flag, bool) {
	if flag, found := lo.Find(c.flags, func(flag *Flag) bool { return flag.name == name }); found {
		return flag, true
	}

	if c.hasParent() {
		return c.parent.findFlag(name)
	}

	return nil, false
}

func dashifyShort(short rune) string {
	return fmt.Sprintf("-%c", short)
}

func isBoolParser(parser argParser) bool {
	_, isBool := parser.Type().(bool)
	return isBool
}
