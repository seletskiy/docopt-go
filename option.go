package docopt

import (
	"bufio"
	"bytes"
)

type Option struct {
	Names       []string
	Default     string
	Value       string
	Description string
	Placeholder string
	IndentLevel int
	LineNumber  int
}

func ParseOptions(section string) ([]Option, error) {
	scanner := bufio.NewScanner(bytes.NewBufferString(section))

	var (
		option  *Option
		options []Option
		index   int
	)

	for scanner.Scan() {
		index += 1

		line := scanner.Text()

		matches, tail := MatcherIndenting.Match(line)

		indenting := matches[1]

		for {
			matches, tail = MatcherOption.Match(tail)
			if matches == nil {
				break
			}

			if option != nil && option.LineNumber != index {
				option = nil
			}

			if option == nil {
				options = append(options, Option{
					IndentLevel: len(indenting),
					LineNumber:  index,
				})

				option = &options[len(options)-1]
			}

			option.Names = append(option.Names, matches[1])

			if matches[2] != "" {
				option.Placeholder = matches[2]
			}

			matches, tail = MatcherDescriptionSeparator.Match(tail)
			if matches != nil {
				break
			}

			matches, tail = MatcherOptionSeparator.Match(tail)
			if matches == nil {
				return nil, ErrParseFailed{
					Message: `expected two or more spaces or option ` +
						`definition, but none found`,

					Line: line,
					Tail: tail,

					LineNumber: index,
				}
			}
		}

		if option != nil {
			if option.LineNumber != index {
				option.Description += " "
			}

			option.Description += tail
		}
	}

	return options, nil
}
