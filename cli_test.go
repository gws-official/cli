package cli

import (
	"context"
	"testing"
)

func TestCommandRunsWithStringFlag(t *testing.T) {
	app := New("tool")

	called := false

	app.Command("greet").
		Flag("name", FlagOptions{
			Alias:    "n",
			Type:     StringFlag,
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			called = true

			if got := c.Flags["name"]; got != "Ada" {
				t.Fatalf("expected name Ada, got %#v", got)
			}

			return nil
		})

	code, err := app.Execute(context.Background(), []string{"greet", "--name", "Ada"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !called {
		t.Fatal("expected command action to be called")
	}
}

func TestAliasFlag(t *testing.T) {
	app := New("tool")

	app.Command("greet").
		Flag("name", FlagOptions{
			Alias:    "n",
			Type:     StringFlag,
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			if got := c.Flags["name"]; got != "Ada" {
				t.Fatalf("expected name Ada, got %#v", got)
			}
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"greet", "-n", "Ada"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestEqualsFlagSyntax(t *testing.T) {
	app := New("tool")

	app.Command("greet").
		Flag("name", FlagOptions{
			Type:     StringFlag,
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			if got := c.Flags["name"]; got != "Ada" {
				t.Fatalf("expected name Ada, got %#v", got)
			}
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"greet", "--name=Ada"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestIntFlag(t *testing.T) {
	app := New("tool")

	app.Command("repeat").
		Flag("times", FlagOptions{
			Type:     IntFlag,
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			if got := c.Flags["times"]; got != 3 {
				t.Fatalf("expected times 3, got %#v", got)
			}
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"repeat", "--times", "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestInvalidIntFlag(t *testing.T) {
	app := New("tool")

	app.Command("repeat").
		Flag("times", FlagOptions{
			Type:     IntFlag,
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"repeat", "--times", "abc"})
	if err == nil {
		t.Fatal("expected error")
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestBoolFlag(t *testing.T) {
	app := New("tool")

	app.Command("build").
		Flag("verbose", FlagOptions{
			Type: BoolFlag,
		}).
		Do(func(ctx context.Context, c *Context) error {
			if got := c.Flags["verbose"]; got != true {
				t.Fatalf("expected verbose true, got %#v", got)
			}
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"build", "--verbose"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestDefaultFlag(t *testing.T) {
	app := New("tool")

	app.Command("repeat").
		Flag("times", FlagOptions{
			Type:    IntFlag,
			Default: 2,
		}).
		Do(func(ctx context.Context, c *Context) error {
			if got := c.Flags["times"]; got != 2 {
				t.Fatalf("expected default times 2, got %#v", got)
			}
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"repeat"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRequiredFlagMissing(t *testing.T) {
	app := New("tool")

	app.Command("greet").
		Flag("name", FlagOptions{
			Type:     StringFlag,
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"greet"})
	if err == nil {
		t.Fatal("expected error")
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestUnknownFlag(t *testing.T) {
	app := New("tool")

	app.Command("greet").
		Do(func(ctx context.Context, c *Context) error {
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"greet", "--wat"})
	if err == nil {
		t.Fatal("expected error")
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestPositionalArg(t *testing.T) {
	app := New("tool")

	app.Command("read").
		Arg("file", ArgOptions{
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			if got := c.Args["file"]; got != "input.txt" {
				t.Fatalf("expected file input.txt, got %#v", got)
			}
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"read", "input.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRequiredArgMissing(t *testing.T) {
	app := New("tool")

	app.Command("read").
		Arg("file", ArgOptions{
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"read"})
	if err == nil {
		t.Fatal("expected error")
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestNestedCommand(t *testing.T) {
	app := New("tool")

	called := false

	app.Command("user").
		Command("add").
		Flag("name", FlagOptions{
			Type:     StringFlag,
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			called = true

			if got := c.Flags["name"]; got != "Ada" {
				t.Fatalf("expected name Ada, got %#v", got)
			}

			if got := c.Command.path(); got != "tool user add" {
				t.Fatalf("expected command path tool user add, got %q", got)
			}

			return nil
		})

	code, err := app.Execute(context.Background(), []string{"user", "add", "--name", "Ada"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !called {
		t.Fatal("expected nested command action to be called")
	}
}

func TestDoubleDashStopsFlagParsing(t *testing.T) {
	app := New("tool")

	app.Command("echo").
		Arg("value", ArgOptions{
			Required: true,
		}).
		Do(func(ctx context.Context, c *Context) error {
			if got := c.Args["value"]; got != "--not-a-flag" {
				t.Fatalf("expected --not-a-flag positional, got %#v", got)
			}
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"echo", "--", "--not-a-flag"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestActionCliErrorCode(t *testing.T) {
	app := New("tool")

	app.Command("fail").
		Do(func(ctx context.Context, c *Context) error {
			return NewError(7, "custom failure")
		})

	code, err := app.Execute(context.Background(), []string{"fail"})
	if err == nil {
		t.Fatal("expected error")
	}
	if code != 7 {
		t.Fatalf("expected exit code 7, got %d", code)
	}
}

func TestHelpDoesNotRunAction(t *testing.T) {
	app := New("tool")

	called := false

	app.Command("greet").
		Description("Greet someone").
		Flag("name", FlagOptions{
			Type: StringFlag,
		}).
		Do(func(ctx context.Context, c *Context) error {
			called = true
			return nil
		})

	code, err := app.Execute(context.Background(), []string{"greet", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if called {
		t.Fatal("expected help not to run action")
	}
}