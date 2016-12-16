package docopt

import (
	"fmt"
	"regexp"
)

type Matcher struct {
	*regexp.Regexp
}

func NewMatcher(expression string, parts ...interface{}) Matcher {
	return Matcher{regexp.MustCompile(`^` + fmt.Sprintf(expression, parts...))}
}

func (matcher *Matcher) Match(body string) ([]string, string) {
	matches := matcher.FindStringSubmatch(body)
	if matches != nil {
		return matches, body[len(matches[0]):]
	}

	return nil, body
}
