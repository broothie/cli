package cli

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/bobg/errors"
	"github.com/samber/lo"
)

//go:embed help.tmpl
var rawHelpTemplate string

var helpTemplate = template.Must(template.New("help").Parse(rawHelpTemplate))

func (c *Command) helpContext() helpContext {
	return helpContext{command: c}
}

func (c *Command) renderHelp(w io.Writer) error {
	if err := helpTemplate.Execute(w, c.helpContext()); err != nil {
		return errors.Wrap(err, "help template")
	}

	return nil
}

type helpContext struct {
	command *Command
}

func (h helpContext) RootName() string {
	return h.command.root().name
}

func (h helpContext) Version() string {
	return h.command.findVersion()
}

func (h helpContext) RootDescription() string {
	return h.command.root().description
}

func (h helpContext) QualifiedName() string {
	return h.command.qualifiedName()
}

func (h helpContext) SubCommands() []*Command {
	return h.command.subCommands
}

func (h helpContext) Arguments() []*Argument {
	return h.command.arguments
}

func (h helpContext) Flags() []*Flag {
	return h.command.flagsUpToRoot()
}

func (h helpContext) SubCommandsTable() (string, error) {
	return tableToString(lo.Map(h.SubCommands(), func(command *Command, _ int) []string {
		return []string{
			"",
			fmt.Sprintf("%s: %s", command.name, command.description),
		}
	}))
}

func (h helpContext) ArgumentList() string {
	return strings.Join(lo.Map(h.Arguments(), func(argument *Argument, _ int) string { return argument.inBrackets() }), " ")
}

func (h helpContext) ArgumentTable() (string, error) {
	return tableToString(lo.Map(h.Arguments(), func(argument *Argument, _ int) []string {
		valueInfo := fmt.Sprintf("(type: %T)", argument.parser.Type())
		if argument.isOptional() {
			valueInfo = fmt.Sprintf("(type: %T, default: %q)", argument.parser.Type(), fmt.Sprint(argument.defaultValue))
		}

		return []string{
			"",
			argument.inBrackets(),
			argument.description,
			valueInfo,
		}
	}))
}

func (h helpContext) FlagTable() (string, error) {
	return tableToString(lo.FilterMap(h.Flags(), func(flag *Flag, _ int) ([]string, bool) {
		if flag.isHidden {
			return nil, false
		}

		longs := lo.Map(append([]string{flag.name}, flag.aliases...), func(long string, _ int) string { return fmt.Sprintf("--%s", long) })

		shorts := ""
		if len(flag.shorts) > 0 {
			shorts = fmt.Sprintf("-%s", string(flag.shorts))
		}

		helpValues := []string{fmt.Sprintf("%q", fmt.Sprint(flag.defaultValue))}
		if flag.defaultEnvName != "" {
			helpValues = append([]string{fmt.Sprintf("$%s", flag.defaultEnvName)}, helpValues...)
		}

		return []string{
			"",
			strings.Join(longs, " "),
			shorts,
			flag.description,
			fmt.Sprintf("(type: %T, default: %s)", flag.parser.Type(), strings.Join(helpValues, ", ")),
		}, true
	}))
}

func tableToString(rows [][]string) (string, error) {
	buffer := new(bytes.Buffer)
	if err := writeTable(buffer, rows); err != nil {
		return "", errors.Wrap(err, "writing table to string")
	}

	return buffer.String(), nil
}

func writeTable(w io.Writer, rows [][]string) error {
	table := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	for _, row := range rows {
		if _, err := fmt.Fprintln(table, strings.Join(row, "\t")); err != nil {
			return errors.Wrap(err, "writing table row")
		}
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "flushing table")
	}

	return nil
}
