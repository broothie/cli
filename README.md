# `cli`

[![Go Reference](https://pkg.go.dev/badge/github.com/broothie/cli.svg)](https://pkg.go.dev/github.com/broothie/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/broothie/cli)](https://goreportcard.com/report/github.com/broothie/cli)
[![codecov](https://codecov.io/gh/broothie/cli/graph/badge.svg?token=GdWFdveewo)](https://codecov.io/gh/broothie/cli)
[![gosec](https://github.com/broothie/cli/actions/workflows/gosec.yml/badge.svg)](https://github.com/broothie/cli/actions/workflows/gosec.yml)

A Go package for building CLIs.

## Installation

```shell
go get github.com/broothie/cli@latest
```

## Documentation

Detailed documentation can be found at [pkg.go.dev](https://pkg.go.dev/github.com/broothie/cli).

## Usage

Using `cli` is as simple as:

```go
// Create and run a command called "fileserver".
// `Run` automatically passes down a `context.Background()` and parses `os.Args[1:]`.
// If an error is returned, and it is either a `cli.ExitError` or an `*exec.ExitError`, the error's exit code will be used.
// For any other errors returned, it exits with code 1.
cli.Run("fileserver", "An HTTP server.",

	// Add an optional positional argument called "root" which will default to ".".
	cli.AddArg("root", "Directory to serve from", cli.SetArgDefault(".")),

	// Add an optional flag called "port" which will default to 3000.
	cli.AddFlag("port", "Port to run server.", cli.SetFlagDefault(3000)),

	// Register a handler for this command.
	// If no handler is registered, it will simply print help and exit.
	cli.SetHandler(func(ctx context.Context) error {
		// Extract the value of the "root" argument.
		root, _ := cli.ArgValue[string](ctx, "root")

		// Extract the value of the "port" flag.
		port, _ := cli.FlagValue[int](ctx, "port")

		addr := fmt.Sprintf(":%d", port)
		return http.ListenAndServe(addr, http.FileServer(http.Dir(root)))
	}),
)

```

## Roadmap

- [ ] Audit bare `err` returns
  - [ ] Two types of errors: config and parse
- [ ] Tab completion
- [ ] Allow variadic arguments
- [ ] Allow slice and map based flags?
