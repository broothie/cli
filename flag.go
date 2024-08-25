package cli

import (
	"context"
	"fmt"
	"reflect"

	"github.com/samber/lo"
)

const helpFlagName = "help"

type Flag struct {
	name         string
	description  string
	aliases      []string
	shorts       []rune
	valueParser  ValueParser
	defaultValue any

	value any
}

func (f *Flag) isHelp() bool {
	return f.name == helpFlagName && f.isBool()
}

func (f *Flag) isBool() bool {
	return isBoolParser(f.valueParser)
}

func FlagValue(ctx context.Context, name string) any {
	cmd := commandFromContext(ctx)
	flag, found := cmd.findFlag(name)
	if !found {
		panic(fmt.Sprintf("no flag with name %q", name))
	}

	if flag.value != nil {
		return flag.value
	}

	return flag.defaultValue
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

func isBoolParser(valueParser ValueParser) bool {
	return reflect.ValueOf(valueParser).Pointer() == reflect.ValueOf(BoolParser).Pointer()
}
