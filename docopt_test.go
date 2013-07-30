package docopt

import (
    "fmt"
    "testing"
    "strings"
)

var _ = fmt.Printf

func TestEmptyUsage(t *testing.T) {
    doc := `Test Empty Usage.
    Usage: prog.go`

    sections := ReSections.FindStringSubmatch(doc)

    if len(sections) != 3 {
        t.Fatal(`can't split to 3 sections`)
    }

    if sections[1] != "prog.go" {
        t.Fatal(`usage section invalid`)
    }

    if sections[2] != "" {
        t.Fatal(`options section invalid`)
    }

    res, err := Docopt(doc)

    if err != nil || res != nil {
        t.Fatal(`invalid result for empty usage`)
    }
}

func TestAltUsage(t *testing.T) {
    usages := `
        prog.go -h
        prog.go --help`
    doc := `Test Alternative Usage.
    Usage:` + usages

    sections := ReSections.FindStringSubmatch(doc)
    if len(sections) != 3 {
        t.Fatal(`can't split to 3 section`)
    }

    if sections[1] != usages {
        t.Fatal(`usage section invalid`)
    }
}

func TestReSections(t *testing.T) {
    usages := `
        prog.go -h
        prog.go --help`
    options := `
        -h --help  Show this message.`
    doc := `Test Alternative Usage.
    Usage:` + usages + `

    Options:` + options

    sections := ReSections.FindStringSubmatch(doc)
    if len(sections) != 3 {
        t.Fatal(`can't split to 3 section`)
    }

    if sections[1] != usages {
        t.Fatal(`usage section invalid`)
    }

    if sections[2] != strings.Trim(options, "\n") {
        t.Fatal(`options section invalid`)
    }

    doc = `Hello, I'm Bogus Test!`
    sections = ReSections.FindStringSubmatch(doc)
    if len(sections) != 0 {
        t.Fatal(`shouldn't match bogus text`)
    }
}

func TestReFlag(t *testing.T) {
    res := ReFlag.FindStringSubmatch(`-h`)
    if res[1] != `-h` {
        t.Fatal(`can't parse -h flag`)
    }

    res = ReFlag.FindStringSubmatch(`-aAFTER`)
    if res[1] != `-a` || res[2] != `AFTER` {
        t.Fatal(`can't parse -aAFTER flag`)
    }

    res = ReFlag.FindStringSubmatch(`-a<after>`)
    if res[1] != `-a` || res[2] != `<after>` {
        t.Fatal(`can't parse -a<after> flag`)
    }

    res = ReFlag.FindStringSubmatch(`-a=AFTER`)
    if res[1] != `-a` || res[2] != `AFTER` {
        t.Fatal(`can't parse -a=AFTER flag`)
    }

    res = ReFlag.FindStringSubmatch(`--help`)
    if res[1] != `--help` {
        t.Fatal(`can't parse --help flag`)
    }

    res = ReFlag.FindStringSubmatch(`--before BEFORE`)
    if res[1] != `--before` || res[2] != `BEFORE` {
        t.Fatal(`can't parse --before BEFORE flag`)
    }

    res = ReFlag.FindStringSubmatch(`--before<before>`)
    if res[1] != `--before` || res[2] != `<before>` {
        t.Fatal(`can't parse --before<before> flag`)
    }

    res = ReFlag.FindStringSubmatch(`--before=<before>`)
    if res[1] != `--before` || res[2] != `<before>` {
        t.Fatal(`can't parse --before=<before> flag`)
    }
}

func TestReOptDesc(t *testing.T) {
    res := ReOptDesc.FindStringSubmatch(`-h  Show help message.`)
    if res[1] != `-h` {
        t.Fatal(`can't parse -h desc`)
    }

    res = ReOptDesc.FindStringSubmatch(`-h --help  Show help message.`)
    if res[1] != `-h --help` {
        t.Fatal(`can't parse -h --help desc`)
    }

    res = ReOptDesc.FindStringSubmatch(`-h, --help  Show help message.`)
    if res[1] != `-h, --help` {
        t.Fatal(`can't parse -h, --help desc`)
    }

    res = ReOptDesc.FindStringSubmatch(`-A AFTER, --after AFTER  Some desc.`)
    if res[1] != `-A AFTER, --after AFTER` {
        t.Fatal(`can't parse -A AFTER, --after AFTER desc`)
    }
}

func TestBogusDoc(t *testing.T) {
    defer func() {
        p := recover()
        if p != `"Usage:" (case-insensetive) not found in help text` {
            t.Fatal(`unexpected panic`)
        }
    }()

    Docopt(`I'm Bogus Help Text!`)

    t.Fatal(`shouldn't reach this point`)
}

func TestHelpOnlyDoc(t *testing.T) {
    doc := `Program only with Help Flags
    Usage: prog.go -h | --help

    Options:
        -h --help  Show this help

    Blah blah blah.
`

    Docopt(doc)
}
