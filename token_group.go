package docopt

import (
	"fmt"
)

type TokenGroup struct {
	Opened   bool
	Required bool

	Pair int
}

func (group *TokenGroup) String() string {
	pairs := []string{"[", "]"}

	if group.Required {
		pairs = []string{"(", ")"}
	}

	pair := pairs[1]
	if group.Opened {
		pair = pairs[0]
	}

	return fmt.Sprintf("%s#%d", pair, group.Pair)
}
