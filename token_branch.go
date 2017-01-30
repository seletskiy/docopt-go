package docopt

import (
	"fmt"
)

type TokenBranch struct {
	Start int
	Next  int
}

func (branch TokenBranch) String() string {
	return fmt.Sprintf(
		"|:%d,%d", branch.Start, branch.Next,
	)
}
