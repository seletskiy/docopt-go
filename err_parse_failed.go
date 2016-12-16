package docopt

import "fmt"

type ErrParseFailed struct {
	Message string
	Line    string
	Tail    string

	LineNumber int
}

func (err ErrParseFailed) Error() string {
	cursor := len(err.Line) - len(err.Tail)
	line := err.Line[:cursor] + "?" + err.Line[cursor:]

	return fmt.Sprintf(
		`%s: (stopped at ?) [L%d] %q`,
		err.Message,
		err.LineNumber,
		line,
	)
}
