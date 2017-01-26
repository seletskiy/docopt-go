package docopt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrammar_Balance_MatchesGroupPairs(t *testing.T) {
	test := assert.New(t)

	grammar := Grammar{
		&TokenStaticWord{},
		&TokenGroup{Opened: true},
		&TokenStaticWord{},
		&TokenGroup{Opened: false},
	}

	err := grammar.Balance()

	test.NoError(err)

	test.Equal(Grammar{
		&TokenStaticWord{},
		&TokenGroup{Opened: true, Pair: 3},
		&TokenStaticWord{},
		&TokenGroup{Opened: false, Pair: 1},
	}, grammar)
}

func TestGrammar_Expand_ExpandsOptionals(t *testing.T) {
	test := assert.New(t)

	grammar := Grammar{
		&TokenStaticWord{Name: "report"},
		&TokenGroup{Opened: true},
		&TokenStaticWord{Name: "verbose"},
		&TokenGroup{Opened: false},
	}

	variants, err := grammar.Expand()

	test.NoError(err)

	fmt.Printf("XXXXXX grammar_test.go:44 variants: %q\n", variants)
}

func TestGrammar_Expand_ExpandsNestedOptionals(t *testing.T) {
	test := assert.New(t)

	grammar := Grammar{
		&TokenStaticWord{Name: "report"},
		&TokenGroup{Opened: true},
		&TokenStaticWord{Name: "verbose"},
		&TokenGroup{Opened: true},
		&TokenStaticWord{Name: "debug"},
		&TokenGroup{Opened: false},
		&TokenGroup{Opened: false},
	}

	variants, err := grammar.Expand()

	test.NoError(err)

	fmt.Printf("XXXXXX grammar_test.go:44 variants: %q\n", variants)
}
