package docopt

import (
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

	test.Equal(
		[]Grammar{
			{
				&TokenStaticWord{Name: "report"},
			},
			{
				&TokenStaticWord{Name: "report"},
				&TokenStaticWord{Name: "verbose"},
			},
		},
		variants,
	)
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

	test.Equal(
		[]Grammar{
			{
				&TokenStaticWord{Name: "report"},
			},
			{
				&TokenStaticWord{Name: "report"},
				&TokenStaticWord{Name: "verbose"},
			},
			{
				&TokenStaticWord{Name: "report"},
				&TokenStaticWord{Name: "verbose"},
				&TokenStaticWord{Name: "debug"},
			},
		},
		variants,
	)
}

func TestGrammar_Expand_ExpandsRequiredGroupWithBranch(t *testing.T) {
	test := assert.New(t)

	grammar := Grammar{
		&TokenGroup{Opened: true, Required: true},
		&TokenStaticWord{Name: "create"},
		&TokenBranch{},
		&TokenStaticWord{Name: "destroy"},
		&TokenGroup{Opened: false, Required: true},
	}

	variants, err := grammar.Expand()

	test.NoError(err)

	test.Equal(
		[]Grammar{
			{
				&TokenStaticWord{Name: "create"},
			},
			{
				&TokenStaticWord{Name: "destroy"},
			},
		},
		variants,
	)
}

func TestGrammar_Expand_ExpandsRequiredGroupWithSeveralBranches(t *testing.T) {
	test := assert.New(t)

	grammar := Grammar{
		&TokenGroup{Opened: true, Required: true},
		&TokenStaticWord{Name: "create"},
		&TokenBranch{},
		&TokenStaticWord{Name: "list"},
		&TokenBranch{},
		&TokenStaticWord{Name: "destroy"},
		&TokenGroup{Opened: false, Required: true},
	}

	variants, err := grammar.Expand()

	test.NoError(err)

	test.Equal(
		[]Grammar{
			{
				&TokenStaticWord{Name: "create"},
			},
			{
				&TokenStaticWord{Name: "list"},
			},
			{
				&TokenStaticWord{Name: "destroy"},
			},
		},
		variants,
	)
}

func TestGrammar_Expand_ExpandsNestedRequiredGroups(t *testing.T) {
	test := assert.New(t)

	grammar := Grammar{
		&TokenGroup{Opened: true, Required: true},
		&TokenStaticWord{Name: "create"},
		&TokenBranch{},
		&TokenStaticWord{Name: "list"},
		&TokenGroup{Opened: true, Required: true},
		&TokenStaticWord{Name: "short"},
		&TokenBranch{},
		&TokenStaticWord{Name: "full"},
		&TokenGroup{Opened: false, Required: true},
		&TokenBranch{},
		&TokenStaticWord{Name: "destroy"},
		&TokenGroup{Opened: false, Required: true},
	}

	variants, err := grammar.Expand()

	test.NoError(err)

	test.Equal(
		[]Grammar{
			{
				&TokenStaticWord{Name: "create"},
			},
			{
				&TokenStaticWord{Name: "list"},
				&TokenStaticWord{Name: "short"},
			},
			{
				&TokenStaticWord{Name: "list"},
				&TokenStaticWord{Name: "full"},
			},
			{
				&TokenStaticWord{Name: "destroy"},
			},
		},
		variants,
	)
}

func TestGrammar_Expand_ExpandsNestedRequiredAndOptionalGroups(t *testing.T) {
	test := assert.New(t)

	grammar := Grammar{
		&TokenGroup{Opened: true, Required: true},
		&TokenStaticWord{Name: "create"},
		&TokenBranch{},
		&TokenStaticWord{Name: "list"},
		&TokenGroup{Opened: true},
		&TokenGroup{Opened: true, Required: true},
		&TokenStaticWord{Name: "short"},
		&TokenBranch{},
		&TokenStaticWord{Name: "full"},
		&TokenGroup{Opened: false, Required: true},
		&TokenGroup{Opened: false},
		&TokenBranch{},
		&TokenStaticWord{Name: "destroy"},
		&TokenGroup{Opened: false, Required: true},
	}

	variants, err := grammar.Expand()

	test.NoError(err)

	test.Equal(
		[]Grammar{
			{
				&TokenStaticWord{Name: "create"},
			},
			{
				&TokenStaticWord{Name: "list"},
			},
			{
				&TokenStaticWord{Name: "list"},
				&TokenStaticWord{Name: "short"},
			},
			{
				&TokenStaticWord{Name: "list"},
				&TokenStaticWord{Name: "full"},
			},
			{
				&TokenStaticWord{Name: "destroy"},
			},
		},
		variants,
	)
}
