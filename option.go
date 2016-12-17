package docopt

import (
	"strings"
)

type Option struct {
	Names       []string
	Value       string
	Description []string
	Placeholder string
	Level       int
}

func (option *Option) GetDescription() string {
	var description string

	for _, line := range option.Description {
		if description == "" {
			description = line
			continue
		}

		if sign, _ := MatcherDescriptionParagraph.Match(line); sign != nil {
			description += "\n" + line
		} else {
			description += " " + line
		}
	}

	return strings.TrimSpace(description)
}

func (option *Option) HasArgument() bool {
	return option.Placeholder != ""
}

func (option *Option) GetDefault() (string, bool) {
	matches, _ := MatcherDescriptionDefault.Match(option.GetDescription())

	if matches == nil {
		return "", false
	}

	return matches[1], true
}
