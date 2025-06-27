package cli_test

import (
	"context"
	"fmt"

	"github.com/broothie/cli"
)

// Example demonstrates using variadic arguments to copy multiple files
func Example_variadic_arguments() {
	cli.Run("copy", "Copy files to destination",
		// Required destination argument  
		cli.AddArg("destination", "Destination directory"),
		
		// Variadic source files argument
		cli.AddArg("sources", "Source files to copy", cli.SetArgVariadic()),
		
		cli.SetHandler(func(ctx context.Context) error {
			// Get the destination
			dest, err := cli.ArgValue[string](ctx, "destination")
			if err != nil {
				return err
			}
			
			// Get all source files as a slice
			sources, err := cli.VariadicArgValue[string](ctx, "sources")
			if err != nil {
				return err
			}
			
			fmt.Printf("Copying %d files to %s:\n", len(sources), dest)
			for _, source := range sources {
				fmt.Printf("  %s -> %s\n", source, dest)
			}
			
			return nil
		}),
	)
}

// Example demonstrates using variadic arguments for a command that processes multiple numbers
func Example_variadic_with_types() {
	cli.Run("sum", "Calculate sum of numbers",
		cli.AddArg("numbers", "Numbers to sum", 
			cli.SetArgParser(cli.IntParser),
			cli.SetArgVariadic(),
		),
		
		cli.SetHandler(func(ctx context.Context) error {
			numbers, err := cli.VariadicArgValue[int](ctx, "numbers")
			if err != nil {
				return err
			}
			
			sum := 0
			for _, num := range numbers {
				sum += num
			}
			
			fmt.Printf("Sum of %v = %d\n", numbers, sum)
			return nil
		}),
	)
}

// Example demonstrates mixing flags with variadic arguments
func Example_variadic_with_flags() {
	cli.Run("process", "Process files with options",
		cli.AddFlag("verbose", "Enable verbose output", cli.SetFlagDefault(false)),
		cli.AddFlag("format", "Output format", cli.SetFlagDefault("text")),
		cli.AddArg("files", "Files to process", cli.SetArgVariadic()),
		
		cli.SetHandler(func(ctx context.Context) error {
			verbose, _ := cli.FlagValue[bool](ctx, "verbose")
			format, _ := cli.FlagValue[string](ctx, "format")
			files, _ := cli.VariadicArgValue[string](ctx, "files")
			
			if verbose {
				fmt.Printf("Processing %d files in %s format:\n", len(files), format)
			}
			
			for _, file := range files {
				if verbose {
					fmt.Printf("Processing: %s\n", file)
				} else {
					fmt.Printf("%s\n", file)
				}
			}
			
			return nil
		}),
	)
}