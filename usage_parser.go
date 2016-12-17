package docopt

type UsageParser struct{}

func (parser *UsageParser) Parse(section string) (*Usage, error) {
	scanner := NewScanner(section)

	usage := &Usage{}

	for scanner.Scan() {
		scanner.Match(MatcherIndenting)

		if scanner.Match(MatcherEndOfLine) != nil {
			continue
		}

		matches := scanner.Match(MatcherUsageWord)
		if matches != nil {
			if usage.Binary == "" {
				usage.Binary = matches[0]

				matches = nil
			}
		}

		scanner.Match(MatcherTokenSeparator)

		grammar := Grammar{}

		for {
			tokens, err := parser.parseTokens(scanner)
			if err != nil {
				return nil, err
			}

			if token == nil {
				break
			}

			grammar = append(grammar, token)
		}

		usage.Variants = append(usage.Variants, grammar)
	}

	return usage, nil
}

func (parser *UsageParser) parseToken(scanner *Scanner) ([]Token, error) {
	if scanner.Match(MatcherEndOfLine) != nil {
		return nil, nil
	}

	tokens := []Token{}

	matches := scanner.Match(MatcherOption)
	if matches != nil {
		tokens = append(tokens, TokenOption{
			Name:        matches[1],
			Placeholder: matches[2],
		})
	}

	matches = scanner.Match(MatcherArgument)
	if matches != nil {
		tokens = append(tokens, TokenPositionalArgument{
			Placeholder: matches[1],
		})
	}

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

	if scanner.Match(MatcherTokenRepeat) != nil {
		tokens = append(tokens, TokenRepeat{})
	}

	scanner.Match(MatcherTokenSeparator)

	if scanner.Match(MatcherTokenBranch) != nil {
		tokens = append(tokens, TokenBranch{})
	}

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

	return nil, scanner.Errorf(
		`expected option, argument, one of "|()[]", "...", ` +
			`but none found`,
	)
}
