package docopt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ArgumentsParser_ParsesArgumentsWithSingleOption(t *testing.T) {
	test := assert.New(t)

	expected := &Arguments{
		Grammar: Grammar{
			TokenOption{Name: "--help"},
		},
	}

	parser := &ArgumentsParser{}

	actual, err := parser.Parse([]string{`--help`})

	test.NoError(err)
	test.EqualValues(expected, actual)
}

func Test_ArgumentsParser_ParsesOptionWithValue(t *testing.T) {
	test := assert.New(t)

	expected := &Arguments{
		Grammar: Grammar{
			TokenOption{Name: "--data", Value: "value"},
		},
	}

	parser := &ArgumentsParser{}

	actual, err := parser.Parse([]string{`--data=value`})

	test.NoError(err)
	test.EqualValues(expected, actual)
}

func Test_ArgumentsParser_ParsesOptionWithSeparatedValue(t *testing.T) {
	test := assert.New(t)

	expected := &Arguments{
		Grammar: Grammar{
			TokenOption{Name: "--data"},
			TokenSeparator{},
			TokenPositionalArgument{Value: `value`},
		},
	}

	parser := &ArgumentsParser{}

	actual, err := parser.Parse([]string{`--data`, `value`})

	test.NoError(err)
	test.EqualValues(expected, actual)
}
