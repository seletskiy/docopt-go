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

func TestParseUsage_ParsesUsageWithRepeat(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah -v...`,
		`blah  -v...`,
		`  blah  -v...`,
		`  blah  -v...  `,
	}

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenOption{Name: "-v"},
				TokenRepeat{},
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

func TestParseUsage_ParsesUsageWithRequiredGroup(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah (--run|--test)`,
		`blah  (--run | --test)`,
		`   blah (--run | --test)`,
		`   blah (--run | --test)   `,
	}

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenGroup{Required: true, Opened: true},
				TokenOption{Name: "--run"},
				TokenBranch{},
				TokenOption{Name: "--test"},
				TokenGroup{Required: true},
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

func TestParseUsage_ParsesUsageWithStaticWord(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah action <value>`,
		`blah  action  <value>`,
		`  blah  action  <value>  `,
	}

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenStaticWord{Name: `action`},
				TokenPositionalArgument{Placeholder: "<value>"},
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

func TestParseUsage_ParsesUsageWithStaticWordsGroup(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah (run|test)`,
		`blah (run | test)`,
		`blah ( run | test )`,
	}

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenGroup{Required: true, Opened: true},
				TokenStaticWord{Name: `run`},
				TokenBranch{},
				TokenStaticWord{Name: `test`},
				TokenGroup{Required: true},
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

func TestParseUsage_ProhibitsEmptyGroup(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah ()`,
		`blah ( )`,
	}

	parser := &UsageParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.Nil(actual)
		test.Error(err)
	}
}

func TestParseUsage_ProhibitsEmptyBranch(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah (|)`,
		`blah ( | )`,
	}

	parser := &UsageParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.Nil(actual)
		test.Error(err)
	}
}

func TestParseUsage_ProhibitsEmptyRepeat(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah ...`,
		`blah (...)`,
		`blah (|...)`,
	}

	parser := &UsageParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.Nil(actual)
		test.Error(err)
	}
}

func TestParseUsage_ParsesUsageWithSeveralUsages(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah --help
		blah --version`,
	}

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{TokenOption{Name: "--help"}},
			{TokenOption{Name: "--version"}},
		},
	}

	parser := &UsageParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.EqualValues(expected, actual)
	}
}

func TestParseUsage_ParsesMultilineUsage(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah --help
			--version`,
	}

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenOption{Name: "--help"},
				TokenOption{Name: "--version"},
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

func TestParseUsage_ProhibitsUsageWithoutBinaryName(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`--help`,
	}

	parser := &UsageParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.Nil(actual)
		test.Error(err)
	}
}

func TestParseUsage_ParsesMultilineUsageStartedWithStaticWord(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`blah action <action>
			value <value>`,
	}

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenStaticWord{Name: "action"},
				TokenPositionalArgument{Placeholder: "<action>"},
				TokenStaticWord{Name: "value"},
				TokenPositionalArgument{Placeholder: "<value>"},
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

func TestParseUsage_ParsesComplexArgumentPatterns(t *testing.T) {
	test := assert.New(t)

	section := `blah action /<search>`

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenStaticWord{Name: "action"},
				TokenPositionalArgument{Placeholder: "/<search>"},
			},
		},
	}

	parser := &UsageParser{}

	actual, err := parser.Parse(section)

	test.NoError(err)
	test.EqualValues(expected, actual)
}
