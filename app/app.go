package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"

	"github.com/paingha/stroll/command"
)

type App struct {
	Name        string
	Description string
	Version     string
	Output      io.Writer
	Context     context.Context
	Args        []string
	cmds        map[string]command.Command
}

var (
	ErrorNameRequired   = errors.New("app name is required")
	ErrorNoArgsProvided = errors.New("no app arguments provided")
)

func (a *App) pre() error {
	if a.Name == "" {
		return ErrorNameRequired
	}
	a.cmds = make(map[string]command.Command, 0)
	switch {
	case a.Args == nil:
		a.Args = os.Args
	case len(a.Args) == 0:
		return ErrorNoArgsProvided
	}
	if a.Output == nil {
		a.Output = os.Stdout
	}
	if a.Context == nil {
		a.Context, _ = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	}
	a.cmds["help"] = command.Command{
		Name:        "help",
		Aka:         "h",
		Description: fmt.Sprintf("shows how to use %s", a.Name),
		Exec: func(ctx context.Context, args []string) error {
			fmt.Fprintf(a.Output, "%s help\n", a.Name)
			return nil
		},
	}
	if _, exists := a.cmds["version"]; !exists {
		a.cmds["version"] = command.Command{
			Name:        "version",
			Aka:         "v",
			Description: fmt.Sprintf("shows version of %s", a.Name),
			Exec: func(ctx context.Context, args []string) error {
				fmt.Fprintf(a.Output, "%s version: %s\n", a.Name, a.Version)
				return nil
			},
		}
	}
	return nil
}

func (a *App) displayCMDs() {
	printCommands(a, a.cmds)
}

func (a *App) Run() error {
	if err := a.pre(); err != nil {
		return err
	}
	if len(os.Args) == 1 {
		a.displayCMDs()
		return nil
	}
	cmd, args := os.Args[1], os.Args[1:]
	if err := a.cmds[cmd].Exec(a.Context, args); err != nil {
		return err
	}
	return nil
}

// PrintCommands in a table form (Name and Description).
func printCommands(cfg *App, cmds map[string]command.Command) {
	minwidth, tabwidth, padding, buffer, flags := 0, 0, 11, byte(' '), uint(0)
	tw := tabwriter.NewWriter(cfg.Output, minwidth, tabwidth, padding, buffer, flags)

	fmt.Fprintf(cfg.Output, "These are common %s commands:\n\n", cfg.Name)

	for _, cmd := range cmds {
		if len(cmd.Children) == 0 {
			printCommand(cfg, tw, "", cmd)
		}

		for _, subcmd := range cmd.Children {
			printCommand(cfg, tw, cmd.Name, subcmd)
		}
	}
	fmt.Fprint(tw, "\n")
	tw.Flush()
}

func printCommand(app *App, writer *tabwriter.Writer, prefix string, cmd command.Command) {
	name := cmd.Name
	if prefix != "" {
		name = fmt.Sprintf("%s %s", prefix, cmd.Name)
	}

	desc := cmd.Description
	if desc == "" {
		desc = "<no description>"
	}

	fmt.Fprintf(writer, "    %s\t%s\n", name, desc)
}
