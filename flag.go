package cli

import (
	"fmt"

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
