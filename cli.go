package cli

import (
	"context"
	"fmt"
	"os"
)

type Action func(ctx context.Context, c *Context) error

type Command struct {
	Name        string
	Desc        string
	Flags       []*Flag
	Args        []*Arg
	Subcommands []*Command
	Action      Action
	Parent      *Command
}

type Context struct {
	Command     *Command
	Flags       map[string]any
	Args        map[string]string
	Positionals []string
	Raw         []string
}

func New(name string) *Command {
	return &Command{Name: name}
}

func (c *Command) Description(desc string) *Command {
	c.Desc = desc
	return c
}

func (c *Command) Command(name string) *Command {
	child := &Command{
		Name:   name,
		Parent: c,
	}
	c.Subcommands = append(c.Subcommands, child)
	return child
}

func (c *Command) Flag(name string, opts FlagOptions) *Command {
	c.Flags = append(c.Flags, &Flag{
		Name:        name,
		Alias:       opts.Alias,
		Type:        opts.Type,
		Description: opts.Description,
		Required:    opts.Required,
		Default:     opts.Default,
		Env:         opts.Env,
	})
	return c
}

func (c *Command) Arg(name string, opts ArgOptions) *Command {
	c.Args = append(c.Args, &Arg{
		Name:        name,
		Description: opts.Description,
		Required:    opts.Required,
		Default:     opts.Default,
	})
	return c
}

func (c *Command) Do(action Action) *Command {
	c.Action = action
	return c
}

func (c *Command) Run(args []string) int {
	code, err := c.Execute(context.Background(), args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	return code
}

func (c *Command) RunAndExit(args []string) {
	os.Exit(c.Run(args))
}

func (c *Command) Execute(ctx context.Context, argv []string) (int, error) {
	cmd, rest := resolveCommand(c, argv)

	if hasHelp(rest) {
		fmt.Println(cmd.Help())
		return 0, nil
	}

	parsed, err := parse(rest, cmd)
	if err != nil {
		return 1, err
	}

	cc, err := buildContext(cmd, argv, parsed)
	if err != nil {
		return 1, err
	}

	if cmd.Action == nil {
		fmt.Println(cmd.Help())
		return 0, nil
	}

	if err := cmd.Action(ctx, cc); err != nil {
		if cliErr, ok := err.(*Error); ok {
			return cliErr.Code, cliErr
		}
		return 1, err
	}

	return 0, nil
}

func resolveCommand(root *Command, argv []string) (*Command, []string) {
	current := root
	remaining := argv

	for len(remaining) > 0 {
		next := current.findSubcommand(remaining[0])
		if next == nil {
			break
		}

		current = next
		remaining = remaining[1:]
	}

	return current, remaining
}

func (c *Command) findSubcommand(name string) *Command {
	for _, sub := range c.Subcommands {
		if sub.Name == name {
			return sub
		}
	}
	return nil
}