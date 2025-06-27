package main

import (
	"context"
	"fmt"
	"os"

	"github.com/broothie/cli"
)

func main() {
	// Create a CLI application with completion support
	cmd, err := cli.NewCommand("fileserver", "An HTTP file server with tab completion",
		// Enable the hidden completion functionality
		cli.EnableCompletion(),
		
		// Add a user-facing completion command  
		cli.AddCompletionCommand(),
		
		// Set version for version flag
		cli.SetVersion("1.0.0"),
		
		// Add version flag
		cli.AddVersionFlag(cli.AddFlagShort('V')),
		
		// Add help flag that's inherited by subcommands
		cli.AddHelpFlag(
			cli.AddFlagShort('h'),
			cli.SetFlagIsInherited(true),
		),
		
		// Add some flags for the main command
		cli.AddFlag("port", "Port to serve on", 
			cli.AddFlagShort('p'),
			cli.SetFlagDefault(8080),
		),
		cli.AddFlag("host", "Host to bind to",
			cli.SetFlagDefault("localhost"),
		),
		cli.AddFlag("verbose", "Enable verbose logging",
			cli.AddFlagShort('v'),
			cli.SetFlagDefault(false),
		),
		
		// Add a positional argument for the directory to serve
		cli.AddArg("directory", "Directory to serve files from",
			cli.SetArgDefault("."),
		),
		
		// Add some subcommands to demonstrate completion
		cli.AddSubCmd("serve", "Start the file server",
			cli.AddFlag("tls", "Enable TLS",
				cli.SetFlagDefault(false),
			),
			cli.AddFlag("cert", "TLS certificate file"),
			cli.AddFlag("key", "TLS private key file"),
			
			cli.SetHandler(func(ctx context.Context) error {
				directory, _ := cli.ArgValue[string](ctx, "directory")
				port, _ := cli.FlagValue[int](ctx, "port")
				host, _ := cli.FlagValue[string](ctx, "host")
				verbose, _ := cli.FlagValue[bool](ctx, "verbose")
				tls, _ := cli.FlagValue[bool](ctx, "tls")
				
				if verbose {
					fmt.Printf("Starting server on %s:%d serving %s\n", host, port, directory)
					if tls {
						fmt.Println("TLS enabled")
					}
				}
				
				fmt.Printf("Server would start on %s:%d serving %s (TLS: %v)\n", host, port, directory, tls)
				return nil
			}),
		),
		
		cli.AddSubCmd("config", "Manage configuration",
			cli.AddSubCmd("show", "Show current configuration",
				cli.SetHandler(func(ctx context.Context) error {
					fmt.Println("Configuration:")
					fmt.Println("  port: 8080")
					fmt.Println("  host: localhost")
					return nil
				}),
			),
			cli.AddSubCmd("set", "Set configuration value",
				cli.AddArg("key", "Configuration key"),
				cli.AddArg("value", "Configuration value"),
				
				cli.SetHandler(func(ctx context.Context) error {
					key, _ := cli.ArgValue[string](ctx, "key")
					value, _ := cli.ArgValue[string](ctx, "value")
					fmt.Printf("Would set %s = %s\n", key, value)
					return nil
				}),
			),
		),
		
		// Default handler for the root command
		cli.SetHandler(func(ctx context.Context) error {
			directory, _ := cli.ArgValue[string](ctx, "directory")
			port, _ := cli.FlagValue[int](ctx, "port")
			host, _ := cli.FlagValue[string](ctx, "host")
			
			fmt.Printf("File server starting on %s:%d serving %s\n", host, port, directory)
			fmt.Println("Use 'fileserver serve' for advanced options")
			fmt.Println("Use 'fileserver completion bash' to generate bash completion")
			return nil
		}),
	)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating command: %v\n", err)
		os.Exit(1)
	}
	
	// Run the command
	if err := cmd.Run(context.Background(), os.Args[1:]); err != nil {
		cli.ExitWithError(err)
	}
}

/*
Example usage and tab completion behavior:

1. Basic command completion:
   $ fileserver <TAB>
   completion  config  serve

2. Flag completion:
   $ fileserver --<TAB>
   --help  --host  --port  --verbose  --version

3. Short flag completion:
   $ fileserver -<TAB>
   -h  -p  -v  -V

4. Subcommand flag completion:
   $ fileserver serve --<TAB>
   --cert  --help  --host  --key  --port  --tls  --verbose

5. Nested subcommand completion:
   $ fileserver config <TAB>
   set  show

6. File/directory completion for arguments:
   $ fileserver /path/to/<TAB>
   [shows files and directories in /path/to/]

7. Generate bash completion script:
   $ fileserver completion bash
   [outputs bash completion script]

To install completion:
   $ fileserver completion bash > /etc/bash_completion.d/fileserver
   $ source /etc/bash_completion.d/fileserver

Or temporarily:
   $ source <(fileserver completion bash)
*/