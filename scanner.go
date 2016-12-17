package docopt

import (
	"bufio"
	"bytes"
	"fmt"
)

type Scanner struct {
	Input *bufio.Scanner

	Line string
	Tail string
}

func NewScanner(input string) *Scanner {
	return &Scanner{
		Input: bufio.NewScanner(bytes.NewBufferString(input)),
	}
}

func (scanner *Scanner) Scan() bool {
	if scanner.Input.Scan() {
		scanner.Line = scanner.Input.Text()
		scanner.Tail = scanner.Line

		return true
	}

	return false
}

func (scanner *Scanner) Match(matcher Matcher) (matches []string) {
	matches, scanner.Tail = matcher.Match(scanner.Tail)

	return matches
}

func (scanner *Scanner) Errorf(message string, args ...interface{}) error {
	return ErrParseFailed{
		Message: fmt.Sprintf(message, args...),
		Line:    scanner.Line,
		Tail:    scanner.Tail,
	}
}
