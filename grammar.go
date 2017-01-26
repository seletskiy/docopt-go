package docopt

import "fmt"

type Grammar []Token

func (grammar Grammar) Balance() error {
	stack := []int{}

	for index, token := range grammar {
		switch token := token.(type) {
		case *TokenGroup:
			if token.Opened {
				stack = append(stack, index)
			} else {
				if len(stack) == 0 {
					return fmt.Errorf(
						`unbalanced group end at position %d`,
						index,
					)
				}

				last := len(stack) - 1

				pair := stack[last]

				grammar[pair].(*TokenGroup).Pair = index
				grammar[index].(*TokenGroup).Pair = pair

				stack[token.Required] = stack[token.Required][:last]
			}
		}
	}

	if len(stack[false])+len(stack[true]) > 0 {
		return fmt.Errorf(
			`one or more of required or optional groups is not closed`,
		)
	}

	return nil
}

func (grammar *Grammar) Expand() ([]Grammar, error) {
	err := grammar.Balance()
	if err != nil {
		return nil, err
	}

	grammars := grammar.expandOptionals()

	for i, grammar := range grammars {
		result := Grammar{}

		for _, token := range grammar {
			if token != nil {
				result = append(result, token)
			}
		}

		grammars[i] = result
	}

	return grammars, nil
}

func (grammar Grammar) expandOptionals() []Grammar {
	stack := []Grammar{
		grammar,
	}

	variants := []Grammar{}

	for len(stack) > 0 {
		found := false

		stack, grammar = stack[:len(stack)-1], stack[len(stack)-1]

	loop:
		for index, token := range grammar {
			switch token := token.(type) {
			case *TokenGroup:
				if token.Required {
					continue
				}

				found = true

				var variant Grammar

				variant = make(Grammar, len(grammar))

				copy(variant, grammar)

				variant[index] = nil
				variant[token.Pair] = nil

				stack = append(stack, variant)

				variant = make(Grammar, len(grammar))

				copy(variant, grammar)

				for i := index; i <= token.Pair; i++ {
					variant[i] = nil
				}

				stack = append(stack, variant)

				break loop
			}
		}

		if !found {
			variants = append(variants, grammar)
		}
	}

	return variants
}
