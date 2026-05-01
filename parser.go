package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type parsedInput struct {
	flags      map[string]string
	boolFlags map[string]bool
	positionals []string
}

func parse(argv []string, cmd *Command) (*parsedInput, error) {
	out := &parsedInput{
		flags:      map[string]string{},
		boolFlags: map[string]bool{},
	}

	for i := 0; i < len(argv); i++ {
		token := argv[i]

		if token == "--" {
			out.positionals = append(out.positionals, argv[i+1:]...)
			break
		}

		if strings.HasPrefix(token, "--") {
			nameValue := strings.TrimPrefix(token, "--")

			if strings.Contains(nameValue, "=") {
				parts := strings.SplitN(nameValue, "=", 2)
				name := parts[0]
				value := parts[1]

				flag := cmd.findFlag(name)
				if flag == nil {
					return nil, fmt.Errorf("unknown flag --%s", name)
				}

				if flag.Type == BoolFlag {
					out.boolFlags[flag.Name] = parseBoolValue(value)
				} else {
					out.flags[flag.Name] = value
				}

				continue
			}

			flag := cmd.findFlag(nameValue)
			if flag == nil {
				return nil, fmt.Errorf("unknown flag --%s", nameValue)
			}

			if flag.Type == BoolFlag {
				out.boolFlags[flag.Name] = true
				continue
			}

			if i+1 >= len(argv) {
				return nil, fmt.Errorf("missing value for --%s", nameValue)
			}

			i++
			out.flags[flag.Name] = argv[i]
			continue
		}

		if strings.HasPrefix(token, "-") && len(token) > 1 {
			alias := strings.TrimPrefix(token, "-")
			flag := cmd.findFlagByAlias(alias)
			if flag == nil {
				return nil, fmt.Errorf("unknown flag -%s", alias)
			}

			if flag.Type == BoolFlag {
				out.boolFlags[flag.Name] = true
				continue
			}

			if i+1 >= len(argv) {
				return nil, fmt.Errorf("missing value for -%s", alias)
			}

			i++
			out.flags[flag.Name] = argv[i]
			continue
		}

		out.positionals = append(out.positionals, token)
	}

	return out, nil
}

func buildContext(cmd *Command, raw []string, parsed *parsedInput) (*Context, error) {
	flags := map[string]any{}

	for _, flag := range cmd.Flags {
		var value any
		var ok bool

		if flag.Type == BoolFlag {
			value, ok = parsed.boolFlags[flag.Name]
		} else {
			value, ok = parsed.flags[flag.Name]
		}

		if !ok && flag.Env != "" {
			if envValue, exists := lookupEnv(flag.Env); exists {
				value = envValue
				ok = true
			}
		}

		if !ok && flag.Default != nil {
			value = flag.Default
			ok = true
		}

		if !ok && flag.Required {
			return nil, fmt.Errorf("missing required flag --%s", flag.Name)
		}

		if ok {
			coerced, err := coerce(value, flag.Type)
			if err != nil {
				return nil, fmt.Errorf("invalid value for --%s: %w", flag.Name, err)
			}
			flags[flag.Name] = coerced
		}
	}

	args := map[string]string{}

	for i, arg := range cmd.Args {
		if i < len(parsed.positionals) {
			args[arg.Name] = parsed.positionals[i]
			continue
		}

		if arg.Default != "" {
			args[arg.Name] = arg.Default
			continue
		}

		if arg.Required {
			return nil, fmt.Errorf("missing required argument %s", arg.Name)
		}
	}

	return &Context{
		Command:     cmd,
		Flags:       flags,
		Args:        args,
		Positionals: parsed.positionals,
		Raw:         raw,
	}, nil
}

func coerce(value any, typ FlagType) (any, error) {
	switch typ {
	case StringFlag:
		return fmt.Sprint(value), nil

	case IntFlag:
		switch v := value.(type) {
		case int:
			return v, nil
		case string:
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			return n, nil
		default:
			return nil, fmt.Errorf("expected int")
		}

	case BoolFlag:
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			return parseBoolValue(v), nil
		default:
			return nil, fmt.Errorf("expected bool")
		}

	default:
		return nil, fmt.Errorf("unknown flag type %s", typ)
	}
}

func parseBoolValue(v string) bool {
	switch strings.ToLower(v) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func (c *Command) findFlag(name string) *Flag {
	for _, f := range c.Flags {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (c *Command) findFlagByAlias(alias string) *Flag {
	for _, f := range c.Flags {
		if f.Alias == alias {
			return f
		}
	}
	return nil
}