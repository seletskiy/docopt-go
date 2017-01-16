package docopt

import "io"

type UsageParser struct{}

func (parser *UsageParser) Parse(section string) (*Usage, error) {
	scanner := NewScanner(section)

	var (
		usage   Usage
		grammar *Grammar
	)

	for scanner.Scan() {
		tokens, err := parser.parseBinaryName(scanner, &usage)
		if err != nil {
			if err == io.EOF {
				continue
			}

			return nil, err
		}

		if usage.Binary == "" {
			return nil, scanner.Errorf(
				`expected binary name`,
			)
		}

		grammar = &usage.Variants[len(usage.Variants)-1]

		*grammar = append(*grammar, tokens...)

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

func (parser *UsageParser) parseBinaryName(
	scanner *Scanner,
	usage *Usage,
) ([]Token, error) {
	scanner.Match(MatcherIndenting)

	if scanner.Match(MatcherEndOfLine) != nil {
		return nil, io.EOF
	}

	matches := scanner.Match(MatcherTokenWord)

	tokens := []Token{
		TokenSeparator{},
	}

	if matches == nil {
		return tokens, nil
	}

	if usage.Binary == "" {
		usage.Binary = matches[0]
	}

	if usage.Binary == matches[0] {
		usage.Variants = append(usage.Variants, Grammar{})

		scanner.Match(MatcherTokenSeparator)

		return nil, nil
	}

	tokens = append(tokens, TokenStaticWord{Name: matches[0]})

	return tokens, nil
}

func (parser *UsageParser) parseTokensGroupStart(
	scanner *Scanner,
) ([]Token, error) {
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

	return tokens, nil
}

func (parser *UsageParser) parseTokensGroupEnd(
	scanner *Scanner,
) ([]Token, error) {
	tokens := []Token{}

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

	return tokens, nil
}

func (parser *UsageParser) parseTokensOptions(
	scanner *Scanner,
) ([]Token, error) {
	tokens := []Token{}

	matches := scanner.Match(MatcherOption)
	if matches != nil {
		tokens = append(tokens, TokenOption{
			Name:  matches[1],
			Value: matches[2],
		})
	}

	matches = scanner.Match(MatcherArgument)
	if matches != nil {
		tokens = append(tokens, TokenPositionalArgument{
			Value: matches[1],
		})
	}

	matches = scanner.Match(MatcherTokenWord)
	if matches != nil {
		tokens = append(tokens, TokenStaticWord{
			Name: matches[0],
		})
	}

	return tokens, nil
}

func (parser *UsageParser) parseTokensRepeat(
	scanner *Scanner,
) ([]Token, error) {
	tokens := []Token{}

	if scanner.Match(MatcherTokenRepeat) != nil {
		tokens = append(tokens, TokenRepeat{})
	}

	return tokens, nil
}

func (parser *UsageParser) parseTokensBranch(
	scanner *Scanner,
) ([]Token, error) {
	tokens := []Token{}

	if scanner.Match(MatcherTokenBranch) != nil {
		tokens = append(tokens, TokenBranch{})

		scanner.Match(MatcherTokenSeparator)
	}

	return tokens, nil
}

func (parser *UsageParser) parseTokens(scanner *Scanner) ([]Token, error) {
	tokens := []Token{}

	if scanner.Match(MatcherTokenSeparator) != nil {
		tokens = append(tokens, TokenSeparator{})
	}

	if scanner.Match(MatcherEndOfLine) != nil {
		return nil, nil
	}

	groups, err := parser.parseTokensGroupStart(scanner)
	if err != nil {
		return nil, err
	}

	empty := false

	if len(groups) > 0 {
		empty = true
	}

	tokens = append(tokens, groups...)

	if scanner.Match(MatcherTokenSeparator) != nil {
		if !empty {
			tokens = append(tokens, TokenSeparator{})
		}
	}

	options, err := parser.parseTokensOptions(scanner)
	if err != nil {
		return nil, err
	}

	if len(options) > 0 {
		empty = false
	}

	tokens = append(tokens, options...)

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

	separator := false
	if scanner.Match(MatcherTokenSeparator) != nil {
		separator = true
	}

	groups, err = parser.parseTokensGroupEnd(scanner)
	if err != nil {
		return nil, err
	}

	if len(groups) > 0 {
		separator = false
	}

	tokens = append(tokens, groups...)

	if len(tokens) == 0 {
		return nil, scanner.Errorf(
			`")" or "]" expected`,
		)
	}

	repeat, err := parser.parseTokensRepeat(scanner)
	if err != nil {
		return nil, err
	}

	if len(repeat) > 0 {
		separator = false
	}

	tokens = append(tokens, repeat...)

	if scanner.Match(MatcherTokenSeparator) != nil {
		separator = true
	}

	branch, err := parser.parseTokensBranch(scanner)
	if err != nil {
		return nil, err
	}

	if len(branch) > 0 {
		separator = false
	}

	tokens = append(tokens, branch...)

	if len(tokens) == 0 {
		return nil, scanner.Errorf(
			`"..." or "|" expected`,
		)
	}

	if separator {
		if scanner.Match(MatcherEndOfLine) == nil {
			tokens = append(tokens, TokenSeparator{})
		}
	}

	return tokens, nil
}
