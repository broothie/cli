package cli

import (
	"context"
	"testing"

	"github.com/broothie/test"
	"github.com/samber/lo"
)

func ensureCalled(t *testing.T) func() {
	called := false
	t.Cleanup(func() { test.True(t, called) })

	return func() { called = true }
}

func Test_git(t *testing.T) {
	type TestCase struct {
		rawArgs         []string
		gitHandler      func(*testing.T) Handler
		commitHandler   func(*testing.T) Handler
		checkoutHandler func(*testing.T) Handler
	}

	testCases := map[string]TestCase{
		"base handler is called": {
			gitHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)

				return func(ctx context.Context) error {
					called()
					return nil
				}
			},
		},
		"long flag is parsed": {
			rawArgs: []string{"--git-dir", "/path/to/something"},
			gitHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)
				return func(ctx context.Context) error {
					called()

					gitDir, err := FlagValue[string](ctx, "git-dir")
					test.NoError(t, err)
					test.Equal(t, "/path/to/something", gitDir)

					return nil
				}
			},
		},
		"long flag is parsed with equal sign": {
			rawArgs: []string{"--git-dir=/path/to/something"},
			gitHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)
				return func(ctx context.Context) error {
					called()

					gitDir, err := FlagValue[string](ctx, "git-dir")
					test.NoError(t, err)
					test.Equal(t, "/path/to/something", gitDir)

					return nil
				}
			},
		},
		"sub-command handler is called": {
			rawArgs: []string{"commit"},
			commitHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)
				return func(ctx context.Context) error {
					called()
					return nil
				}
			},
		},
		"flag alias is parsed": {
			rawArgs: []string{"commit", "--msg", "a commit message"},
			commitHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)
				return func(ctx context.Context) error {
					called()

					message, err := FlagValue[string](ctx, "message")
					test.NoError(t, err)
					test.Equal(t, "a commit message", message)

					return nil
				}
			},
		},
		"flag short is parsed": {
			rawArgs: []string{"commit", "-m", "a commit message"},
			commitHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)
				return func(ctx context.Context) error {
					called()

					isAll, err := FlagValue[bool](ctx, "all")
					test.NoError(t, err)
					test.False(t, isAll)

					message, err := FlagValue[string](ctx, "message")
					test.NoError(t, err)
					test.Equal(t, "a commit message", message)

					return nil
				}
			},
		},
		"flag short with equal sign is parsed": {
			rawArgs: []string{"commit", "-m=a commit message"},
			commitHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)
				return func(ctx context.Context) error {
					called()

					isAll, err := FlagValue[bool](ctx, "all")
					test.NoError(t, err)
					test.False(t, isAll)

					message, err := FlagValue[string](ctx, "message")
					test.NoError(t, err)
					test.Equal(t, "a commit message", message)

					return nil
				}
			},
		},
		"short flag group is parsed": {
			rawArgs: []string{"commit", "-am", "a commit message"},
			commitHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)
				return func(ctx context.Context) error {
					called()

					isAll, err := FlagValue[bool](ctx, "all")
					test.NoError(t, err)
					test.True(t, isAll)

					message, err := FlagValue[string](ctx, "message")
					test.NoError(t, err)
					test.Equal(t, "a commit message", message)

					return nil
				}
			},
		},
		"arg is parsed": {
			rawArgs: []string{"checkout", "some-branch"},
			checkoutHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)
				return func(ctx context.Context) error {
					called()

					branch, err := ArgValue[string](ctx, "branch")
					test.NoError(t, err)
					test.Equal(t, "some-branch", branch)

					return nil
				}
			},
		},
		"bool flag is parsed with arg": {
			rawArgs: []string{"checkout", "-b", "some-branch"},
			checkoutHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)
				return func(ctx context.Context) error {
					called()

					branch, err := ArgValue[string](ctx, "branch")
					test.NoError(t, err)
					test.Equal(t, "some-branch", branch)

					isNewBranch, err := FlagValue[bool](ctx, "new-branch")
					test.NoError(t, err)
					test.True(t, isNewBranch)

					return nil
				}
			},
		},
		"env var based flag is evaluated": {
			gitHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)

				return func(ctx context.Context) error {
					called()

					globalGitignore, err := FlagValue[string](ctx, "global-gitignore")
					test.NoError(t, err)
					test.Equal(t, globalGitignore, "path/to/some/.gitignore")
					return nil
				}
			},
		},
		"rest args": {
			rawArgs: []string{"checkout", "some-branch", "--", "more", "tokens", "here"},
			checkoutHandler: func(t *testing.T) Handler {
				called := ensureCalled(t)

				return func(ctx context.Context) error {
					called()

					rest, err := Rest(ctx)
					test.NoError(t, err)
					test.DeepEqual(t, rest, []string{"more", "tokens", "here"})
					return nil
				}
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("GLOBAL_GITIGNORE", "path/to/some/.gitignore")

			command, err := NewCommand("git", "the stupid content tracker",
				SetVersion("v0.1.0"),
				AddFlag("git-dir", "Git directory to use"),
				AddFlag("global-gitignore", "Global .gitignore file to use.",
					SetFlagDefaultEnv("GLOBAL_GITIGNORE"),
				),
				AddSubCmd("commit", "Record changes to the repository",
					AddFlag("message", "commit message",
						AddFlagAlias("msg"),
						AddFlagShort('m'),
					),
					AddFlag("all", "commit all changed files",
						AddFlagShort('a'),
						SetFlagDefault(false),
					),
					SetHandler(lo.IfF(testCase.commitHandler != nil, func() Handler { return testCase.commitHandler(t) }).Else(nil)),
				),
				AddSubCmd("checkout", "Switch branches or restore working tree files",
					AddArg("branch", "Branch to check out"),
					AddFlag("new-branch", "New branch name",
						AddFlagShort('b'),
						SetFlagDefault(false),
					),
					SetHandler(lo.IfF(testCase.checkoutHandler != nil, func() Handler { return testCase.checkoutHandler(t) }).Else(nil)),
				),
				SetHandler(lo.IfF(testCase.gitHandler != nil, func() Handler { return testCase.gitHandler(t) }).Else(nil)),
			)

			test.NoError(t, err)
			test.Nil(t, command.Run(context.TODO(), testCase.rawArgs))
		})
	}
}
