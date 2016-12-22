package docopt

import (
	"fmt"
	"regexp"
)

type Matcher struct {
	*regexp.Regexp
}

func NewMatcher(expression string, parts ...interface{}) Matcher {
	return Matcher{regexp.MustCompile(fmt.Sprintf(expression, parts...))}
}

func (matcher *Matcher) Match(body string) ([]string, string) {
	index := matcher.FindStringIndex(body)
	if index != nil {
		if index[0] != 0 {
			return nil, body
		}

		return matcher.FindStringSubmatch(body), body[index[1]:]
	}

	return nil, body
}
