package docopt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OptionsParser_ParsesSingleOptionWithDescription(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`-a  An single option.`,
		`-a   An single option.`,
	}

	expected := []Option{
		{
			Names:       []string{"-a"},
			Description: []string{`An single option.`},
		},
	}

	parser := OptionsParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.EqualValues(expected, actual)
	}
}

func Test_OptionsParser_ParsesShortAndLongOptionWithDescription(t *testing.T) {
	test := assert.New(t)

	section := `-a --an-option  An interesting option.`

	expected := []Option{
		{
			Names:       []string{"-a", "--an-option"},
			Description: []string{`An interesting option.`},
		},
	}

	parser := OptionsParser{}

	actual, err := parser.Parse(section)

	test.NoError(err)
	test.EqualValues(expected, actual)
}

func Test_OptionsParser_JoinsDescriptionOnSeveralLines(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`-a  An single
		     option.`,

		`-a  An single
		      option.`,

		`-a  An
		      single
		      option.`,
	}

	parser := OptionsParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.Len(actual, 1)
		test.EqualValues(actual[0].GetDescription(), `An single option.`)
	}
}

func Test_OptionsParser_ParsesDescriptionOnNewLine(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`-a
			An single option.`,
		`-a
			An single
			option.`,
	}

	parser := OptionsParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.Len(actual, 1)
		test.EqualValues(actual[0].GetDescription(), `An single option.`)
	}
}

func Test_OptionsParser_ParsesValueForShortOption(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`-a <value>  Option with value.`,
		`-a=<value>  Option with value.`,
		`-a<value>   Option with value.`,
		`-a VALUE    Option with value.`,
		`-a=VALUE    Option with value.`,
		`-aVALUE     Option with value.`,
	}

	parser := OptionsParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.Len(actual, 1)
		test.True(actual[0].HasArgument())
	}
}

func Test_OptionsParser_ParsesValueForOptionWithAliases(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`-a --an-option <value>  Option with value.`,
		`-a --an-option=<value>  Option with value.`,
		`-a --an-option<value>   Option with value.`,
		`-a --an-option VALUE    Option with value.`,
		`-a --an-option=VALUE    Option with value.`,
	}

	parser := OptionsParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.Len(actual, 1)
		test.True(actual[0].HasArgument())
	}
}

func Test_OptionsParser_ParsesOptionDefaultValue(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`-a <value>  Option with value [default: some].`,
		`-a <value>  Option with value
					 [default: some].`,
	}

	parser := OptionsParser{}

	for _, variant := range variants {
		actual, err := parser.Parse(variant)

		test.NoError(err)
		test.Len(actual, 1)
		test.True(actual[0].HasArgument())

		_default, has := actual[0].GetDefault()
		test.True(has)
		test.Equal(`some`, _default)
	}
}

func Test_OptionsParser_ParsesSeveralOptions(t *testing.T) {
	test := assert.New(t)

	section := `
		-a --an-option  First option.
		-o              Another
		                 option.

		-x --with-value <value>  Option with value.
		                         [default: xxx]
	`

	parser := OptionsParser{}

	options, err := parser.Parse(section)
	test.NoError(err)
	test.Len(options, 3)

	test.False(options[0].HasArgument())
	test.EqualValues([]string{`-a`, `--an-option`}, options[0].Names)
	test.EqualValues(`First option.`, options[0].GetDescription())

	test.False(options[1].HasArgument())
	test.EqualValues([]string{`-o`}, options[1].Names)
	test.EqualValues(`Another option.`, options[1].GetDescription())

	test.True(options[2].HasArgument())
	test.EqualValues([]string{`-x`, `--with-value`}, options[2].Names)
	test.EqualValues(`<value>`, options[2].Value)
}
