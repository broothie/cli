package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samber/lo"
)

// CompletionResult represents a single completion suggestion
type CompletionResult struct {
	Value       string
	Description string
}

// CompletionContext contains information about the current completion request
type CompletionContext struct {
	Command     *Command
	Args        []string
	CurrentWord string
	PreviousWord string
	WordIndex   int
}

// Completer is a function that generates completion suggestions
type Completer func(ctx CompletionContext) []CompletionResult

// FileCompleter generates file and directory completions
func FileCompleter(ctx CompletionContext) []CompletionResult {
	var results []CompletionResult
	
	pattern := ctx.CurrentWord
	if pattern == "" {
		pattern = "*"
	}
	
	// Handle glob patterns
	if !strings.Contains(pattern, "*") && !strings.Contains(pattern, "?") {
		pattern = pattern + "*"
	}
	
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return results
	}
	
	for _, match := range matches {
		stat, err := os.Stat(match)
		if err != nil {
			continue
		}
		
		result := CompletionResult{Value: match}
		if stat.IsDir() {
			result.Description = "directory"
			result.Value = match + "/"
		} else {
			result.Description = "file"
		}
		
		results = append(results, result)
	}
	
	return results
}

// DirectoryCompleter generates only directory completions
func DirectoryCompleter(ctx CompletionContext) []CompletionResult {
	var results []CompletionResult
	
	pattern := ctx.CurrentWord
	if pattern == "" {
		pattern = "*"
	}
	
	if !strings.Contains(pattern, "*") && !strings.Contains(pattern, "?") {
		pattern = pattern + "*"
	}
	
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return results
	}
	
	for _, match := range matches {
		stat, err := os.Stat(match)
		if err != nil || !stat.IsDir() {
			continue
		}
		
		results = append(results, CompletionResult{
			Value:       match + "/",
			Description: "directory",
		})
	}
	
	return results
}

// generateCompletions generates completion results for the given context
func (c *Command) generateCompletions(ctx CompletionContext) []CompletionResult {
	var results []CompletionResult
	
	// If current word starts with a flag prefix, complete flags
	if strings.HasPrefix(ctx.CurrentWord, "--") {
		results = append(results, c.completeLongFlags(ctx)...)
	} else if strings.HasPrefix(ctx.CurrentWord, "-") && len(ctx.CurrentWord) > 1 {
		results = append(results, c.completeShortFlags(ctx)...)
	} else {
		// Check if we're expecting a flag value
		if ctx.PreviousWord != "" && c.isFlagExpectingValue(ctx.PreviousWord) {
			// For now, provide file completion for flag values
			results = append(results, FileCompleter(ctx)...)
		} else {
			// Complete subcommands first
			results = append(results, c.completeSubcommands(ctx)...)
			
			// If no subcommands match, try arguments
			if len(results) == 0 {
				results = append(results, c.completeArguments(ctx)...)
			}
		}
	}
	
	return results
}

// completeLongFlags generates completions for long flags (--flag)
func (c *Command) completeLongFlags(ctx CompletionContext) []CompletionResult {
	var results []CompletionResult
	prefix := strings.TrimPrefix(ctx.CurrentWord, "--")
	
	flags := c.flagsUpToRoot()
	for _, flag := range flags {
		if flag.isHidden {
			continue
		}
		
		// Check flag name
		if strings.HasPrefix(flag.name, prefix) {
			results = append(results, CompletionResult{
				Value:       "--" + flag.name,
				Description: flag.description,
			})
		}
		
		// Check flag aliases
		for _, alias := range flag.aliases {
			if strings.HasPrefix(alias, prefix) {
				results = append(results, CompletionResult{
					Value:       "--" + alias,
					Description: flag.description + " (alias)",
				})
			}
		}
	}
	
	return results
}

// completeShortFlags generates completions for short flags (-f)
func (c *Command) completeShortFlags(ctx CompletionContext) []CompletionResult {
	var results []CompletionResult
	
	// For short flags, we only complete single character flags
	if len(ctx.CurrentWord) != 2 {
		return results
	}
	
	prefix := ctx.CurrentWord[1] // Get the character after '-'
	
	flags := c.flagsUpToRoot()
	for _, flag := range flags {
		if flag.isHidden {
			continue
		}
		
		for _, short := range flag.shorts {
			if rune(prefix) == short {
				results = append(results, CompletionResult{
					Value:       fmt.Sprintf("-%c", short),
					Description: flag.description,
				})
			}
		}
	}
	
	return results
}

// completeSubcommands generates completions for subcommands
func (c *Command) completeSubcommands(ctx CompletionContext) []CompletionResult {
	var results []CompletionResult
	
	for _, subCmd := range c.subCommands {
		if strings.HasPrefix(subCmd.name, ctx.CurrentWord) {
			results = append(results, CompletionResult{
				Value:       subCmd.name,
				Description: subCmd.description,
			})
		}
		
		// Check aliases
		for _, alias := range subCmd.aliases {
			if strings.HasPrefix(alias, ctx.CurrentWord) {
				results = append(results, CompletionResult{
					Value:       alias,
					Description: subCmd.description + " (alias)",
				})
			}
		}
	}
	
	return results
}

// completeArguments generates completions for positional arguments
func (c *Command) completeArguments(ctx CompletionContext) []CompletionResult {
	// Count non-flag arguments to determine which argument we're completing
	argCount := c.countNonFlagArgs(ctx.Args[:ctx.WordIndex])
	
	if argCount >= len(c.arguments) {
		return []CompletionResult{}
	}
	
	// For now, provide file completion for all arguments
	// This can be enhanced later with custom argument completers
	return FileCompleter(ctx)
}

// countNonFlagArgs counts arguments that are not flags or flag values
func (c *Command) countNonFlagArgs(args []string) int {
	count := 0
	skipNext := false
	
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}
		
		if strings.HasPrefix(arg, "-") {
			// Check if this flag expects a value
			if c.isFlagExpectingValue(arg) && i+1 < len(args) {
				skipNext = true
			}
		} else {
			// Check if this is a subcommand
			if _, found := lo.Find(c.subCommands, func(cmd *Command) bool { 
				return cmd.name == arg || lo.Contains(cmd.aliases, arg)
			}); !found {
				count++
			}
		}
	}
	
	return count
}

// isFlagExpectingValue checks if a flag expects a value
func (c *Command) isFlagExpectingValue(flagArg string) bool {
	var flag *Flag
	var found bool
	
	if strings.HasPrefix(flagArg, "--") {
		flagName := strings.TrimPrefix(flagArg, "--")
		flag, found = c.findLongFlag(flagName)
	} else if strings.HasPrefix(flagArg, "-") && len(flagArg) == 2 {
		short := rune(flagArg[1])
		flag, found = c.findShortFlag(short)
	}
	
	if !found {
		return false
	}
	
	return !flag.isBool()
}

// Complete generates completion suggestions for the given arguments
func (c *Command) Complete(args []string) []CompletionResult {
	if len(args) == 0 {
		return c.generateCompletions(CompletionContext{
			Command:     c,
			Args:        args,
			CurrentWord: "",
			WordIndex:   0,
		})
	}
	
	// Find the current word being completed (last argument)
	currentWord := args[len(args)-1]
	previousWord := ""
	if len(args) > 1 {
		previousWord = args[len(args)-2]
	}
	
	ctx := CompletionContext{
		Command:      c,
		Args:         args,
		CurrentWord:  currentWord,
		PreviousWord: previousWord,
		WordIndex:    len(args) - 1,
	}
	
	// Check if we need to delegate to a subcommand
	for _, subCmd := range c.subCommands {
		if len(args) > 0 && (subCmd.name == args[0] || lo.Contains(subCmd.aliases, args[0])) {
			// Delegate to subcommand
			return subCmd.Complete(args[1:])
		}
	}
	
	results := c.generateCompletions(ctx)
	
	// Sort results alphabetically
	sort.Slice(results, func(i, j int) bool {
		return results[i].Value < results[j].Value
	})
	
	return results
}

// GenerateBashCompletion generates a bash completion script for the command
func (c *Command) GenerateBashCompletion() string {
	rootName := c.root().name
	
	script := fmt.Sprintf(`#!/bin/bash

_%s_completion() {
    local cur prev words cword
    _init_completion || return

    # Set up completion environment
    export COMP_LINE="${COMP_LINE}"
    export COMP_POINT="${COMP_POINT}"

    # Get completions from the command
    local completions
    completions=$(%s __complete 2>/dev/null)

    if [[ $? -eq 0 ]]; then
        COMPREPLY=($(compgen -W "${completions}" -- "${cur}"))
    fi
}

# Register the completion function
complete -F _%s_completion %s
`, rootName, rootName, rootName, rootName)

	return script
}

// CompletionCommand creates a hidden completion command
func CompletionCommand(rootCmd *Command) *Command {
	cmd, _ := NewCommand("__complete", "Generate shell completions (hidden)",
		SetHandler(func(ctx context.Context) error {
			// Parse completion arguments from environment or command line
			compLine := os.Getenv("COMP_LINE")
			if compLine == "" {
				// Fallback: try to get from remaining args
				if len(os.Args) > 2 {
					compLine = strings.Join(os.Args[2:], " ")
				}
			}
			
			if compLine == "" {
				return nil
			}
			
			// Parse the completion line
			args := strings.Fields(compLine)
			if len(args) > 0 {
				// Remove the first argument (the command name itself)
				args = args[1:]
			}
			
			completions := rootCmd.Complete(args)
			for _, completion := range completions {
				fmt.Println(completion.Value)
			}
			
			return nil
		}),
	)
	
	return cmd
}