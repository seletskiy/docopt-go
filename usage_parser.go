package docopt

type UsageParser struct{}

func (parser *UsageParser) Parse(section string) (*Usage, error) {
	scanner := NewScanner(section)

	var (
		usage   Usage
		grammar *Grammar
	)

	for scanner.Scan() {
		scanner.Match(MatcherIndenting)

		if scanner.Match(MatcherEndOfLine) != nil {
			continue
		}

		matches := scanner.Match(MatcherTokenWord)
		if matches != nil {
			if usage.Binary == "" {
				usage.Binary = matches[0]
			}

			if usage.Binary == matches[0] {
				usage.Variants = append(usage.Variants, Grammar{})

				grammar = &usage.Variants[len(usage.Variants)-1]
			} else {
				*grammar = append(*grammar, TokenStaticWord{
					Name: matches[0],
				})
			}
		}

		if usage.Binary == "" {
			return nil, scanner.Errorf(
				`expected binary name`,
			)
		}

		for {
			tokens, err := parser.parseTokens(scanner)
			if err != nil {
				return nil, err
			}

			if tokens == nil {
				break
			}

			*grammar = append(*grammar, tokens...)
		}
	}

	return &usage, nil
}

func (parser *UsageParser) parseTokens(scanner *Scanner) ([]Token, error) {
	scanner.Match(MatcherTokenSeparator)

	if scanner.Match(MatcherEndOfLine) != nil {
		return nil, nil
	}

	tokens := []Token{}

	if scanner.Match(MatcherTokenRequiredGroupStart) != nil {
		tokens = append(tokens, TokenGroup{
			Opened:   true,
			Required: true,
		})
	}

	if scanner.Match(MatcherTokenOptionalGroupStart) != nil {
		tokens = append(tokens, TokenGroup{
			Opened:   true,
			Required: false,
		})
	}

	empty := false
	if len(tokens) > 0 {
		empty = true
	}

	scanner.Match(MatcherTokenSeparator)

	matches := scanner.Match(MatcherOption)
	if matches != nil {
		empty = false

		tokens = append(tokens, TokenOption{
			Name:        matches[1],
			Placeholder: matches[2],
		})
	}

	matches = scanner.Match(MatcherArgument)
	if matches != nil {
		empty = false

		tokens = append(tokens, TokenPositionalArgument{
			Placeholder: matches[1],
		})
	}

	matches = scanner.Match(MatcherTokenWord)
	if matches != nil {
		empty = false

		tokens = append(tokens, TokenStaticWord{
			Name: matches[0],
		})
	}

	if len(tokens) == 0 {
		return nil, scanner.Errorf(
			`static word, option or argument expected`,
		)
	}

	if empty {
		return nil, scanner.Errorf(
			`empty groups not allowed`,
		)
	}

	scanner.Match(MatcherTokenSeparator)

	if scanner.Match(MatcherTokenRequiredGroupEnd) != nil {
		tokens = append(tokens, TokenGroup{
			Opened:   false,
			Required: true,
		})
	}

	if scanner.Match(MatcherTokenOptionalGroupEnd) != nil {
		tokens = append(tokens, TokenGroup{
			Opened:   false,
			Required: false,
		})
	}

	if len(tokens) == 0 {
		return nil, scanner.Errorf(
			`")" or "]" expected`,
		)
	}

	if scanner.Match(MatcherTokenRepeat) != nil {
		tokens = append(tokens, TokenRepeat{})
	}

	scanner.Match(MatcherTokenSeparator)

	if scanner.Match(MatcherTokenBranch) != nil {
		tokens = append(tokens, TokenBranch{})
	}

	if len(tokens) == 0 {
		return nil, scanner.Errorf(
			`"..." or "|" expected`,
		)
	}

	return tokens, nil
}
