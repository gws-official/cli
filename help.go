package cli

import (
	"fmt"
	"strings"
)

func hasHelp(argv []string) bool {
	for _, arg := range argv {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func (c *Command) Help() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Usage:\n  %s", c.path())

	if len(c.Subcommands) > 0 {
		b.WriteString(" <command>")
	}

	if len(c.Flags) > 0 {
		b.WriteString(" [options]")
	}

	for _, arg := range c.Args {
		if arg.Required {
			fmt.Fprintf(&b, " <%s>", arg.Name)
		} else {
			fmt.Fprintf(&b, " [%s]", arg.Name)
		}
	}

	b.WriteString("\n")

	if c.Desc != "" {
		fmt.Fprintf(&b, "\nDescription:\n  %s\n", c.Desc)
	}

	if len(c.Subcommands) > 0 {
		b.WriteString("\nCommands:\n")
		for _, sub := range c.Subcommands {
			fmt.Fprintf(&b, "  %-16s %s\n", sub.Name, sub.Desc)
		}
	}

	if len(c.Args) > 0 {
		b.WriteString("\nArguments:\n")
		for _, arg := range c.Args {
			fmt.Fprintf(&b, "  %-16s %s\n", arg.Name, arg.Description)
		}
	}

	b.WriteString("\nOptions:\n")
	for _, flag := range c.Flags {
		name := "--" + flag.Name
		if flag.Alias != "" {
			name = "-" + flag.Alias + ", " + name
		}

		if flag.Type != BoolFlag {
			name += " <" + string(flag.Type) + ">"
		}

		fmt.Fprintf(&b, "  %-24s %s\n", name, flag.Description)
	}

	b.WriteString("  -h, --help               Show help\n")

	return b.String()
}

func (c *Command) path() string {
	var parts []string

	for current := c; current != nil; current = current.Parent {
		parts = append([]string{current.Name}, parts...)
	}

	return strings.Join(parts, " ")
}