package docopt

import "fmt"

type Grammar []Token

func (grammar Grammar) Balance() error {
	pairs := []int{}

	for index, token := range grammar {
		var (
			group, _  = token.(*TokenGroup)
			branch, _ = token.(*TokenBranch)
		)

		if group == nil && branch == nil {
			continue
		}

		if group != nil && group.Opened {
			pairs = append(pairs, index)

			group.Pair = index

			continue
		}

		if len(pairs) == 0 {
			return fmt.Errorf(
				"unbalanced group end at position %d",
				index,
			)
		}

		var (
			pair  = pairs[len(pairs)-1]
			start = grammar[pair].(*TokenGroup)
		)

		if branch != nil {
			branch.Start = pair
		}

		if group != nil {
			if start.Required != group.Required {
				return fmt.Errorf(
					"opened group at position %d do not match %d",
					pair, index,
				)
			}

			group.Pair = pair

			pairs = pairs[:len(pairs)-1]
		}

		if branch, ok := grammar[start.Pair].(*TokenBranch); ok {
			branch.Next = index
		}

		start.Pair = index
	}

	if len(pairs) > 0 {
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

	grammars := grammar.expandGroups()

	for i, grammar := range grammars {
		normalized := Grammar{}

		for _, token := range grammar {
			if token != nil {
				normalized = append(normalized, token)
			}
		}

		grammars[i] = normalized
	}

	return grammars, nil
}

func (grammar Grammar) expandGroups() []Grammar {
	stack := []Grammar{grammar}

	variants := []Grammar{}

	for len(stack) > 0 {
		found := false

		stack, grammar = stack[:len(stack)-1], stack[len(stack)-1]

		for index, token := range grammar {
			if token, ok := token.(*TokenGroup); ok {
				found = true

				if !token.Required {
					stack = append(stack, grammar.cut(index, token.Pair))
				}

				stack = append(
					stack,
					grammar.cut(index, index).cut(token.Pair, token.Pair),
				)

				break
			}

			if token, ok := token.(*TokenBranch); ok {
				found = true

				end := token.Next

				for {
					token, ok := grammar[end].(*TokenBranch)
					if !ok {
						break
					}

					end = token.Next
				}

				stack = append(
					stack,
					grammar.cut(index, end),
				)

				for {
					token, ok := grammar[index].(*TokenBranch)
					if !ok {
						break
					}

					stack = append(
						stack,
						grammar.cut(token.Start, index).cut(token.Next, end),
					)

					index = token.Next
				}

				break
			}
		}

		if !found {
			variants = append([]Grammar{grammar}, variants...)
		}
	}

	return variants
}

func (grammar Grammar) cut(begin, end int) Grammar {
	var clone = make(Grammar, len(grammar))

	copy(clone, grammar)

	for i := begin; i <= end; i++ {
		clone[i] = nil
	}

	return clone
}
