package docopt

import "fmt"

const (
	CursorSign = `→`
)

type ErrParseFailed struct {
	Message string
	Line    string
	Tail    string
}

func (err ErrParseFailed) Error() string {
	cursor := len(err.Line) - len(err.Tail)
	line := err.Line[:cursor] + CursorSign + err.Line[cursor:]

	return fmt.Sprintf(`%s: %q`, err.Message, line)
}
