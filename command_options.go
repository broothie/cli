package cli

import (
	"github.com/broothie/option"
)

func SetVersion(version string) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		command.version = version
		return command, nil
	}
}

func AddAlias(alias string) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		command.aliases = append(command.aliases, alias)
		return command, nil
	}
}

func SetHandler(handler Handler) option.Func[*Command] {
	return func(command *Command) (*Command, error) {
		command.handler = handler
		return command, nil
	}
}

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

func AddHelpFlag(options ...option.Option[*Flag]) option.Func[*Command] {
	defaultOptions := option.NewOptions(SetFlagDefault(false))
	return AddFlag(helpFlagName, "Print help.", append(defaultOptions, options...)...)
}
