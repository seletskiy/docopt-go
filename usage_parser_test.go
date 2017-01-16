package docopt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UsageParser_ParsesEmptyUsage(t *testing.T) {
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

func Test_UsageParser_ParsesUsageWithSingleOption(t *testing.T) {
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

func Test_UsageParser_ParsesUsageWithSingleBranch(t *testing.T) {
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

func Test_UsageParser_ParsesUsageWithRepeat(t *testing.T) {
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

func Test_UsageParser_ParsesUsageWithRequiredGroup(t *testing.T) {
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

func Test_UsageParser_ParsesUsageWithStaticWord(t *testing.T) {
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
				TokenSeparator{},
				TokenPositionalArgument{Value: "<value>"},
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

func Test_UsageParser_ParsesUsageWithStaticWordsGroup(t *testing.T) {
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

func Test_UsageParser_ProhibitsEmptyGroup(t *testing.T) {
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

func Test_UsageParser_ProhibitsEmptyBranch(t *testing.T) {
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

func Test_UsageParser_ProhibitsEmptyRepeat(t *testing.T) {
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

func Test_UsageParser_ParsesUsageWithSeveralUsages(t *testing.T) {
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

func Test_UsageParser_ParsesMultilineUsage(t *testing.T) {
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
				TokenSeparator{},
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

func Test_UsageParser_ProhibitsUsageWithoutBinaryName(t *testing.T) {
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

func Test_UsageParser_ParsesMultilineUsageStartedWithStaticWord(t *testing.T) {
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
				TokenSeparator{},
				TokenPositionalArgument{Value: "<action>"},
				TokenSeparator{},
				TokenStaticWord{Name: "value"},
				TokenSeparator{},
				TokenPositionalArgument{Value: "<value>"},
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

func Test_UsageParser_ParsesComplexArgumentPatterns(t *testing.T) {
	test := assert.New(t)

	section := `blah action /<search>`

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenStaticWord{Name: "action"},
				TokenSeparator{},
				TokenStaticWord{Name: "/"},
				TokenPositionalArgument{Value: "<search>"},
			},
		},
	}

	parser := &UsageParser{}

	actual, err := parser.Parse(section)

	test.NoError(err)
	test.EqualValues(expected, actual)
}

func Test_UsageParser_ParsesComplexArgumentPatternsWithGroups(t *testing.T) {
	test := assert.New(t)

	section := `blah action key[=<value>]`

	expected := &Usage{
		Binary: "blah",
		Variants: []Grammar{
			{
				TokenStaticWord{Name: "action"},
				TokenSeparator{},
				TokenStaticWord{Name: "key"},
				TokenGroup{Opened: true},
				TokenStaticWord{Name: "="},
				TokenPositionalArgument{Value: "<value>"},
				TokenGroup{Opened: false},
			},
		},
	}

	parser := &UsageParser{}

	actual, err := parser.Parse(section)

	test.NoError(err)
	test.EqualValues(expected, actual)
}
