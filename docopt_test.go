package docopt

import "testing"

func TestEmptyUsage(t *testing.T) {
    res, err := Docopt(`Test Empty Usage.
        Usage: prog.go`)

    if err != nil || res != nil {
        t.Fail()
    }
}
