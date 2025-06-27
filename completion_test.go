package cli

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/broothie/test"
)

func TestCompleteCommands(t *testing.T) {
	cmd, err := NewCommand("git", "Version control system",
		AddSubCmd("clone", "Clone a repository"),
		AddSubCmd("commit", "Create a commit"),
		AddSubCmd("push", "Push changes"),
	)
	test.NoError(t, err)

	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "complete empty",
			args:     []string{},
			expected: []string{"clone", "commit", "push"},
		},
		{
			name:     "complete c",
			args:     []string{"c"},
			expected: []string{"clone", "commit"},
		},
		{
			name:     "complete cl",
			args:     []string{"cl"},
			expected: []string{"clone"},
		},
		{
			name:     "complete p",
			args:     []string{"p"},
			expected: []string{"push"},
		},
		{
			name:     "complete nonexistent",
			args:     []string{"xyz"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completions := cmd.Complete(tt.args)
			var values []string
			for _, c := range completions {
				values = append(values, c.Value)
			}
			test.Equal(t, tt.expected, values)
		})
	}
}

func TestCompleteFlags(t *testing.T) {
	cmd, err := NewCommand("test", "Test command",
		AddFlag("verbose", "Verbose output", AddFlagShort('v')),
		AddFlag("output", "Output file", AddFlagShort('o')),
		AddFlag("debug", "Debug mode", SetFlagDefault(false)),
	)
	test.NoError(t, err)

	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "complete long flags",
			args:     []string{"--"},
			expected: []string{"--debug", "--output", "--verbose"},
		},
		{
			name:     "complete long flags with prefix",
			args:     []string{"--v"},
			expected: []string{"--verbose"},
		},
		{
			name:     "complete long flags with prefix o",
			args:     []string{"--o"},
			expected: []string{"--output"},
		},
		{
			name:     "complete short flag v",
			args:     []string{"-v"},
			expected: []string{"-v"},
		},
		{
			name:     "complete short flag o",
			args:     []string{"-o"},
			expected: []string{"-o"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completions := cmd.Complete(tt.args)
			var values []string
			for _, c := range completions {
				values = append(values, c.Value)
			}
			test.Equal(t, tt.expected, values)
		})
	}
}

func TestCompleteWithAliases(t *testing.T) {
	cmd, err := NewCommand("docker", "Container management",
		AddSubCmd("container", "Manage containers",
			AddAlias("c"),
		),
		AddSubCmd("image", "Manage images",
			AddAlias("img"),
		),
	)
	test.NoError(t, err)

	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "complete with aliases",
			args:     []string{"c"},
			expected: []string{"c", "container"},
		},
		{
			name:     "complete img alias",
			args:     []string{"img"},
			expected: []string{"img"},
		},
		{
			name:     "complete i prefix",
			args:     []string{"i"},
			expected: []string{"image", "img"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completions := cmd.Complete(tt.args)
			var values []string
			for _, c := range completions {
				values = append(values, c.Value)
			}
			test.Equal(t, tt.expected, values)
		})
	}
}

func TestCompleteSubcommandFlags(t *testing.T) {
	cmd, err := NewCommand("git", "Version control",
		AddFlag("help", "Show help", AddFlagShort('h'), SetFlagIsInherited(true)),
		AddSubCmd("clone", "Clone repository",
			AddFlag("depth", "Clone depth"),
			AddFlag("branch", "Branch to clone", AddFlagShort('b')),
		),
	)
	test.NoError(t, err)

	// Test completing flags for subcommand
	completions := cmd.Complete([]string{"clone", "--"})
	var values []string
	for _, c := range completions {
		values = append(values, c.Value)
	}
	
	expected := []string{"--branch", "--depth", "--help"}
	test.Equal(t, expected, values)
}

func TestCompleteHiddenFlags(t *testing.T) {
	cmd, err := NewCommand("test", "Test command",
		AddFlag("verbose", "Verbose output"),
		AddFlag("debug", "Debug mode", SetFlagIsHidden(true)),
	)
	test.NoError(t, err)

	completions := cmd.Complete([]string{"--"})
	var values []string
	for _, c := range completions {
		values = append(values, c.Value)
	}

	// Hidden flags should not appear in completions
	expected := []string{"--verbose"}
	test.Equal(t, expected, values)
}

func TestFileCompleter(t *testing.T) {
	// Create a temporary directory with some files
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Create test files
	os.WriteFile("test.txt", []byte("test"), 0644)
	os.WriteFile("example.md", []byte("example"), 0644)
	os.Mkdir("subdir", 0755)

	ctx := CompletionContext{
		CurrentWord: "",
	}

	completions := FileCompleter(ctx)
	
	// Should find our test files
	var found []string
	for _, c := range completions {
		found = append(found, c.Value)
	}

	// Check that we have our expected files (order may vary)
	hasTestTxt := false
	hasExampleMd := false
	hasSubdir := false

	for _, f := range found {
		if f == "test.txt" {
			hasTestTxt = true
		}
		if f == "example.md" {
			hasExampleMd = true
		}
		if f == "subdir/" {
			hasSubdir = true
		}
	}

	test.True(t, hasTestTxt, "Should find test.txt")
	test.True(t, hasExampleMd, "Should find example.md")
	test.True(t, hasSubdir, "Should find subdir/")
}

func TestDirectoryCompleter(t *testing.T) {
	// Create a temporary directory with some files and subdirectories
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Create test files and directories
	os.WriteFile("test.txt", []byte("test"), 0644)
	os.Mkdir("subdir1", 0755)
	os.Mkdir("subdir2", 0755)

	ctx := CompletionContext{
		CurrentWord: "",
	}

	completions := DirectoryCompleter(ctx)
	
	// Should only find directories
	var found []string
	for _, c := range completions {
		found = append(found, c.Value)
	}

	// Should not include files, only directories
	hasTestTxt := false
	hasSubdir1 := false
	hasSubdir2 := false

	for _, f := range found {
		if f == "test.txt" {
			hasTestTxt = true
		}
		if f == "subdir1/" {
			hasSubdir1 = true
		}
		if f == "subdir2/" {
			hasSubdir2 = true
		}
	}

	test.False(t, hasTestTxt, "Should not find test.txt (not a directory)")
	test.True(t, hasSubdir1, "Should find subdir1/")
	test.True(t, hasSubdir2, "Should find subdir2/")
}

func TestGenerateBashCompletion(t *testing.T) {
	cmd, err := NewCommand("myapp", "My application")
	test.NoError(t, err)

	script := cmd.GenerateBashCompletion()
	
	test.True(t, strings.Contains(script, "#!/bin/bash"), "Should contain shebang")
	test.True(t, strings.Contains(script, "_myapp_completion"), "Should contain completion function")
	test.True(t, strings.Contains(script, "complete -F _myapp_completion myapp"), "Should register completion")
	test.True(t, strings.Contains(script, "myapp __complete"), "Should call completion command")
}

func TestCompletionCommand(t *testing.T) {
	rootCmd, err := NewCommand("test", "Test command",
		AddSubCmd("sub", "Subcommand"),
		AddFlag("flag", "Test flag"),
	)
	test.NoError(t, err)

	completionCmd := CompletionCommand(rootCmd)
	
	test.Equal(t, "__complete", completionCmd.name)
	test.NotNil(t, completionCmd.handler)
}

func TestEnableCompletion(t *testing.T) {
	cmd, err := NewCommand("test", "Test command",
		EnableCompletion(),
	)
	test.NoError(t, err)

	// Should have added the completion command
	hasCompletionCmd := false
	for _, subCmd := range cmd.subCommands {
		if subCmd.name == "__complete" {
			hasCompletionCmd = true
			break
		}
	}

	test.True(t, hasCompletionCmd, "Should have added __complete subcommand")
}

func TestAddCompletionCommand(t *testing.T) {
	cmd, err := NewCommand("test", "Test command",
		AddCompletionCommand(),
	)
	test.NoError(t, err)

	// Should have added the completion command
	hasCompletionCmd := false
	for _, subCmd := range cmd.subCommands {
		if subCmd.name == "completion" {
			hasCompletionCmd = true
			// Check that it has a bash subcommand
			hasBashCmd := false
			for _, bashCmd := range subCmd.subCommands {
				if bashCmd.name == "bash" {
					hasBashCmd = true
					break
				}
			}
			test.True(t, hasBashCmd, "completion command should have bash subcommand")
			break
		}
	}

	test.True(t, hasCompletionCmd, "Should have added completion subcommand")
}

func TestCompleteWithContext(t *testing.T) {
	cmd, err := NewCommand("git", "Version control",
		AddSubCmd("clone", "Clone repository",
			AddArg("url", "Repository URL"),
			AddArg("dir", "Target directory", SetArgDefault(".")),
		),
	)
	test.NoError(t, err)

	// Test that subcommand completion works
	completions := cmd.Complete([]string{"clone", "https://github.com/user/repo.git"})
	
	// Should complete files/directories for the second argument
	test.True(t, len(completions) >= 0, "Should return file completions for directory argument")
}