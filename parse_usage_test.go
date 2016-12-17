package docopt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUsage_ParsesEmptyUsage(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah`,
		`blah  `,
		`  blah  `,
	}

	expected := &Usage{
		Binary:   "blah",
		Variants: []Grammar{Grammar{}},
	}

	parser := &UsageParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.EqualValues(expected, actual)
	}
}

func TestParseUsage_ParsesUsageWithSingleOption(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah --help`,
		`blah  --help`,
		`  blah  --help`,
	}

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{TokenOption{Name: "--help"}},
		},
	}

	parser := &UsageParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.EqualValues(expected, actual)
	}
}

func TestParseUsage_ParsesUsageWithSingleBranch(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah --help|-h`,
		`blah --help | -h`,
		`blah  --help | -h`,
		`  blah  --help | -h`,
	}

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenOption{Name: "--help"},
				TokenBranch{},
				TokenOption{Name: "-h"},
			},
		},
	}

	parser := &UsageParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.EqualValues(expected, actual)
	}
}
