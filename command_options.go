package cli

import "github.com/broothie/option"

const versionFlagName = "version"

// SetVersion sets the version of the command.
func SetVersion(version string) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		command.version = version
		return command, nil
	}
}

// AddAlias adds an alias to the command.
func AddAlias(alias string) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		command.aliases = append(command.aliases, alias)
		return command, nil
	}
}

func SetInstallZshTabCompletion(installZshTabCompletion bool) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		command.installZshTabCompletion = installZshTabCompletion
		return command, nil
	}
}

// SetHandler sets the handler of the command.
func SetHandler(handler Handler) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		command.handler = handler
		return command, nil
	}
}

// AddSubCmd adds a subcommand to the command.
func AddSubCmd(name, description string, options ...option.Option[*Command]) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		subCommand, err := NewCommand(name, description, options...)
		if err != nil {
			return nil, err
		}

		subCommand.parent = command
		command.subCommands = append(command.subCommands, subCommand)
		return command, nil
	}
}

// AddFlag adds a flag to the command.
func AddFlag(name, description string, options ...option.Option[*Flag]) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		flag, err := newFlag(name, description, options...)
		if err != nil {
			return nil, err
		}

		command.flags = append(command.flags, flag)
		return command, nil
	}
}

// AddArg adds an argument to the command.
func AddArg(name, description string, options ...option.Option[*Argument]) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		argument, err := newArgument(name, description, options...)
		if err != nil {
			return nil, err
		}

		command.arguments = append(command.arguments, argument)
		return command, nil
	}
}

// AddHelpFlag adds a help flag to the command.
func AddHelpFlag(options ...option.Option[*Flag]) option.Func[*Command] {
	defaultOptions := option.NewOptions(SetFlagDefault(false))
	return AddFlag(helpFlagName, "Print help.", append(defaultOptions, options...)...)
}

func AddVersionFlag(options ...option.Option[*Flag]) option.Func[*Command] {
	defaultOptions := option.NewOptions(SetFlagDefault(false))
	return AddFlag(versionFlagName, "Print version.", append(defaultOptions, options...)...)
}
