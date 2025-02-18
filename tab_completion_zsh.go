package cli

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bobg/errors"
	"github.com/samber/lo"
)

var (
	//go:embed tab-completion.zsh.tmpl
	rawZshTabCompletionTemplate string

	zshTabCompletionTemplate = template.Must(template.New("tab-completion.zsh").Parse(rawZshTabCompletionTemplate))
)

func (c *Command) installZshCompletion() error {
	zshPath, err := exec.LookPath("zsh")
	if errors.Is(err, exec.ErrNotFound) {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "finding zsh executable")
	}

	// Check if completion is already installed
	out, err := exec.Command(zshPath, "-c", "echo $fpath").Output()
	if err != nil {
		return errors.Wrap(err, "getting zsh fpath")
	}

	// Parse space-separated paths
	paths := strings.Fields(string(out))
	if len(paths) == 0 {
		return errors.New("no completion paths found in $fpath")
	}

	for _, directory := range paths {
		directory := strings.Trim(directory, "()")
		completionPath := filepath.Join(directory, fmt.Sprintf("_%s", c.name))

		// Check if already installed
		if _, err := os.Stat(completionPath); err == nil {
			return nil
		}

		if err := c.createAndWriteCompletion(completionPath); err != nil {
			return err
		}
	}

	return errors.New("no writable paths found in $fpath")
}

func (c *Command) createAndWriteCompletion(completionPath string) (err error) {
	f, err := os.Create(completionPath)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	return c.renderZshAutocompleteScript(f)
}

func (c *Command) renderZshAutocompleteScript(w io.Writer) error {
	if err := zshTabCompletionTemplate.Execute(w, c.zshAutocompleteContext()); err != nil {
		return errors.Wrap(err, "zsh autocomplete template")
	}

	return nil
}

func (c *Command) zshAutocompleteContext() map[string]any {
	return map[string]any{
		"name": c.name,
		"subCommands": lo.Map(c.subCommands, func(command *Command, _ int) map[string]any {
			return map[string]any{
				"name":        command.name,
				"description": command.description,
			}
		}),
		"flags": lo.Map(c.flags, func(flag *Flag, _ int) map[string]any {
			return map[string]any{
				"name":        flag.name,
				"description": flag.description,
				"aliases":     flag.aliases,
				"shorts":      flag.shorts,
			}
		}),
	}
}
