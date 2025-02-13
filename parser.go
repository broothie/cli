package cli

import (
	"context"
	"strings"

	"github.com/bobg/errors"
	"github.com/samber/lo"
)

const (
	flagPrefix     = "-"
	longFlagPrefix = "--"
)

var (
	InvalidFlagError        = errors.New("invalid flag")
	MissingFlagValueError   = errors.New("missing flag value")
	TooManyArgumentsError   = errors.New("too many arguments")
	FlagGroupWithEqualError = errors.New("short flags with equal signs cannot be grouped")
)

type parser struct {
	command *Command
	tokens  []string

	index         int
	argumentIndex int
	errors        []error
}

func newParser(command *Command, tokens []string) *parser {
	return &parser{
		command: command,
		tokens:  tokens,
	}
}

func (c *Command) newParser(tokens []string) *parser {
	return newParser(c, tokens)
}

func (p *parser) parse(ctx context.Context) error {
	for p.index < len(p.tokens) {
		commandProcessed, err := p.parseArg(ctx)
		if err != nil {
			return err
		} else if commandProcessed {
			return nil
		}
	}

	if err := p.command.validateInput(); err != nil {
		return err
	}

	return p.command.runHandler(ctx)
}

func (p *parser) parseArg(ctx context.Context) (bool, error) {
	current, _ := p.current()

	if strings.HasPrefix(current, flagPrefix) {
		return false, p.processFlag()
	} else if command, found := lo.Find(p.command.subCommands, func(subCommand *Command) bool { return subCommand.name == current }); found {
		return true, p.processCommand(ctx, command)
	}

	return false, p.processArg()
}

func (p *parser) processFlag() error {
	current, _ := p.current()

	if strings.HasPrefix(current, longFlagPrefix) {
		return p.processLongFlag()
	}

	return p.processShortFlagGroup()
}

func (p *parser) processLongFlag() error {
	current, _ := p.current()

	if strings.Contains(current, "=") {
		return p.processLongFlagWithEqual()
	}

	flag, found := p.command.findLongFlag(strings.TrimPrefix(current, longFlagPrefix))
	if !found {
		return errors.Wrapf(InvalidFlagError, "no flag found for %q", current)
	}

	if flag.isBool() {
		if flag.isHelp() {
			p.command.handler = helpHandler
		}

		flag.value = !flag.defaultValue.(bool)
		p.index += 1
		return nil
	}

	next, nextPresent := p.next()
	if !nextPresent {
		return errors.Wrapf(MissingFlagValueError, "flag %q", current)
	}

	value, err := flag.parser.Parse(next)
	if err != nil {
		return errors.Wrapf(err, "parsing provided value %q for flag %q", next, current)
	}

	flag.value = value
	p.index += 2
	return nil
}

func (p *parser) processLongFlagWithEqual() error {
	current, _ := p.current()

	rawFlag, rawValue, _ := strings.Cut(current, "=")
	flag, found := p.command.findLongFlag(strings.TrimPrefix(rawFlag, longFlagPrefix))
	if !found {
		return errors.Wrapf(InvalidFlagError, "no flag found for %q", rawFlag)
	}

	if flag.isHelp() {
		p.command.handler = helpHandler
	}

	value, err := flag.parser.Parse(rawValue)
	if err != nil {
		return errors.Wrapf(err, "parsing provided value %q for flag %q", rawValue, rawFlag)
	}

	flag.value = value
	p.index += 1
	return nil
}

func (p *parser) processShortFlagGroup() error {
	current, _ := p.current()
	if strings.Contains(current, "=") {
		return p.processShortFlagWithEqual()
	}

	incrementIndexBy := 1
	for _, short := range strings.TrimPrefix(current, flagPrefix) {
		wasValueProcessed, err := p.processShortFlag(short)
		if err != nil {
			return err
		}

		if wasValueProcessed {
			incrementIndexBy = 2
		}
	}

	p.index += incrementIndexBy
	return nil
}

func (p *parser) processShortFlag(short rune) (bool, error) {
	flag, found := p.command.findShortFlag(short)
	if !found {
		return false, errors.Wrapf(InvalidFlagError, "no short flag found for %q", dashifyShort(short))
	}

	if flag.isBool() {
		if flag.isHelp() {
			p.command.handler = helpHandler
		}

		flag.value = !flag.defaultValue.(bool)
		return false, nil
	}

	next, nextPresent := p.next()
	if !nextPresent {
		return false, errors.Wrapf(MissingFlagValueError, "flag %q", dashifyShort(short))
	}

	value, err := flag.parser.Parse(next)
	if err != nil {
		return false, errors.Wrapf(err, "parsing provided value %q for flag %q", next, dashifyShort(short))
	}

	flag.value = value
	return true, nil
}

func (p *parser) processShortFlagWithEqual() error {
	current, _ := p.current()

	rawFlag, rawValue, _ := strings.Cut(current, "=")
	flagName := strings.TrimPrefix(rawFlag, flagPrefix)
	if len(flagName) != 1 {
		return errors.Wrapf(FlagGroupWithEqualError, "flag %q", current)
	}

	short := rune(flagName[0])
	flag, found := p.command.findShortFlag(short)
	if !found {
		return errors.Wrapf(InvalidFlagError, "no short flag found for %q", dashifyShort(short))
	}

	if flag.isHelp() {
		p.command.handler = helpHandler
	}

	value, err := flag.parser.Parse(rawValue)
	if err != nil {
		return errors.Wrapf(err, "parsing provided value %q for flag %q", rawValue, dashifyShort(short))
	}

	flag.value = value
	p.index += 1
	return nil
}

func (p *parser) processCommand(ctx context.Context, command *Command) error {
	return command.Run(ctx, p.unprocessed())
}

func (p *parser) processArg() error {
	if p.argumentIndex >= len(p.command.arguments) {
		return errors.Wrapf(TooManyArgumentsError, "only expected %d arguments", len(p.command.arguments))
	}

	current, _ := p.current()
	argument := p.command.arguments[p.argumentIndex]
	value, err := argument.parser.Parse(current)
	if err != nil {
		return errors.Wrapf(err, "parsing provided value %q for argument %d", current, p.argumentIndex+1)
	}

	argument.value = value
	p.index += 1
	p.argumentIndex += 1
	return nil
}

func (c *Command) findLongFlag(name string) (*Flag, bool) {
	flag, found := lo.Find(c.flags, func(flag *Flag) bool { return flag.name == name || lo.Contains(flag.aliases, name) })
	if found {
		return flag, true
	}

	if c.hasParent() {
		return c.parent.findLongFlag(name)
	}

	return nil, false
}

func (c *Command) findShortFlag(short rune) (*Flag, bool) {
	flag, found := lo.Find(c.flags, func(flag *Flag) bool { return lo.Contains(flag.shorts, short) })
	if found {
		return flag, true
	}

	if c.hasParent() {
		return c.parent.findShortFlag(short)
	}

	return nil, false
}

func (p *parser) current() (string, bool) {
	return p.atOffset(0)
}

func (p *parser) next() (string, bool) {
	return p.atOffset(1)
}

func (p *parser) atOffset(offset int) (string, bool) {
	return p.at(p.index + offset)
}

func (p *parser) at(index int) (string, bool) {
	if !p.indexIsInBounds(index) {
		return "", false
	}

	return p.tokens[index], true
}

func (p *parser) indexIsInBounds(index int) bool {
	return 0 <= index && index < len(p.tokens)
}

func (p *parser) unprocessed() []string {
	return p.tokens[p.index+1:]
}
