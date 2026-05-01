package cli

import "fmt"

type FlagType string

const (
	StringFlag FlagType = "string"
	IntFlag    FlagType = "int"
	BoolFlag   FlagType = "bool"
)

type FlagOptions struct {
	Alias       string
	Type        FlagType
	Description string
	Required    bool
	Default     any
	Env         string
}

type ArgOptions struct {
	Description string
	Required    bool
	Default     string
}

type Flag struct {
	Name        string
	Alias       string
	Type        FlagType
	Description string
	Required    bool
	Default     any
	Env         string
}

type Arg struct {
	Name        string
	Description string
	Required    bool
	Default     string
}

type Error struct {
	Message string
	Code    int
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}