package docopt

type ArgumentsParser struct{}

func (parser *ArgumentsParser) Parse(args []string) (*Arguments, error) {
	var (
		arguments Arguments
	)

	for _, arg := range args {
		if len(arguments.Grammar) > 0 {
			arguments.Grammar = append(arguments.Grammar, TokenSeparator{})
		}

		scanner := NewScanner(arg)

		for scanner.Scan() {
			token := parser.parseToken(scanner)

			if token == nil {
				break
			}

			arguments.Grammar = append(arguments.Grammar, token)
		}
	}

	return &arguments, nil
}

func (parser *ArgumentsParser) parseTokenOption(
	scanner *Scanner,
) Token {
	matches := scanner.Match(MatcherOptionName)
	if matches != nil {
		token := TokenOption{
			Name: matches[1],
		}

		scanner.Match(MatcherOptionValueSeparator)

		matches = scanner.Match(MatcherAny)
		if matches != nil {
			token.Value = matches[0]
		}

		return token
	}

	return nil
}

func (parser *ArgumentsParser) parseTokenValue(
	scanner *Scanner,
) Token {
	matches := scanner.Match(MatcherAny)
	if matches != nil {
		return TokenPositionalArgument{
			Value: matches[0],
		}
	}

	return nil
}

func (parser *ArgumentsParser) parseToken(scanner *Scanner) Token {
	option := parser.parseTokenOption(scanner)
	if option != nil {
		return option
	}

	value := parser.parseTokenValue(scanner)
	if value != nil {
		return value
	}

	return nil
}
