package docopt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions_ParsesSingleOptionWithDescription(t *testing.T) {
	test := assert.New(t)

	variants := []string{
		`-a  An single option.`,
		`-a   An single option.`,
		`-a   An single option.`,
		`-a   An single
			   option.`,
	}

	expected := []Option{
		{
			Names:       []string{"-a"},
			Description: `An single option.`,
			LineNumber:  1,
		},
	}

	for _, variant := range variants {
		actual, err := ParseOptions(variant)

		test.NoError(err)
		test.EqualValues(expected, actual)
	}
}

func TestParseOptions_ParsesShortAndLongOptionWithDescription(t *testing.T) {
	test := assert.New(t)

	section := `-a --an-option  An interesting option.`

	expected := []Option{
		{
			Names:       []string{"-a", "--an-option"},
			Description: `An interesting option.`,
			LineNumber:  1,
		},
	}

	actual, err := ParseOptions(section)
	test.NoError(err)
	test.EqualValues(expected, actual)
}
