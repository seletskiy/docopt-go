package docopt

type OptionsParser struct{}

func (parser *OptionsParser) Parse(section string) ([]Option, error) {
	scanner := NewScanner(section)

	var (
		option  *Option
		options []Option
	)

	for scanner.Scan() {
		matches := scanner.Match(MatcherIndenting)

		indenting := matches[1]

		allocate := true

		for {
			matches = scanner.Match(MatcherOption)
			if matches == nil {
				break
			}

			if allocate {
				option = nil
			}

			if option == nil {
				options = append(options, Option{
					Level: len(indenting),
				})

				option = &options[len(options)-1]

				allocate = false
			}

			option.Names = append(option.Names, matches[1])

			if matches[2] != "" {
				option.Value = matches[2]
			}

			if scanner.Match(MatcherDescriptionSeparator) != nil {
				break
			}

			if scanner.Match(MatcherEndOfLine) != nil {
				break
			}

			matches = scanner.Match(MatcherOptionSeparator)
			if matches == nil {
				return nil, scanner.Errorf(
					`expected two or more spaces or option ` +
						`definition, but none found`,
				)
			}
		}

		if option != nil {
			option.Description = append(option.Description, scanner.Tail)
		}
	}

	return options, nil
}
