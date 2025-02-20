package cli

import (
	"bytes"
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
	"golang.org/x/sync/errgroup"
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

		if err := c.createAndWriteCompletion(completionPath); err == nil {
			return nil
		}
	}

	return errors.New("no writable paths found in $fpath")
}

func (c *Command) createAndWriteCompletion(completionPath string) (err error) {
	file, err := os.OpenFile(completionPath, os.O_RDWR|os.O_CREATE, 0666)
	defer func() {
		if file == nil {
			return
		}

		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	if os.IsNotExist(err) {
		if err := c.renderZshAutocompleteScript(file); err != nil {
			return errors.Wrap(err, "writing tab completion script")
		}
	} else if err != nil {
		return errors.Wrap(err, "opening completion script")
	}

	group := new(errgroup.Group)

	var fileContents []byte
	group.Go(func() error {
		var err error
		if fileContents, err = io.ReadAll(file); err != nil {
			return errors.Wrap(err, "reading existing tab completion script")
		}

		return nil
	})

	buffer := new(bytes.Buffer)
	group.Go(func() error {
		if err := c.renderZshAutocompleteScript(buffer); err != nil {
			return errors.Wrap(err, "generating tab completion script")
		}

		return nil
	})

	if err := group.Wait(); err != nil {
		return err
	}

	if bytes.Equal(buffer.Bytes(), fileContents) {
		return nil
	}

	if err := file.Truncate(0); err != nil {
		return errors.Wrap(err, "truncating existing tab completion script")
	}

	fmt.Println("writing script", completionPath)
	if _, err := io.Copy(file, buffer); err != nil {
		return errors.Wrap(err, "writing tab completion script")
	}

	return nil
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
		"args": lo.Map(c.arguments, func(argument *Argument, _ int) map[string]any {
			return map[string]any{
				"name":        argument.name,
				"description": argument.description,
			}
		}),
		"flags": lo.Map(c.flags, func(flag *Flag, _ int) map[string]any {
			return map[string]any{
				"name":        flag.name,
				"description": flag.description,
				"aliases":     flag.aliases,
				"shorts":      string(flag.shorts),
			}
		}),
	}
}
