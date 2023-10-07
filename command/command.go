package command

import (
	"context"
	"flag"
)

type Command struct {
	// Name is the name of the CLI command.
	Name string
	// Aka is the alias name for the CLI command.
	Aka string
	// Description is the help explanation of what the command does.
	Description string
	// Children are subcommands of this command.
	Children []Command
	// Exec is called when the command is called.
	Exec func(ctx context.Context, args []string) error
	// Flags is the list of flags associated with the command.
	Flags flag.FlagSet
}

func (c *Command) run(ctx context.Context, args []string) error {
	return c.Exec(ctx, args)
}
