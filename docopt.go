package docopt

import (
    "fmt"
    re "regexp"
    "strings"
)

var _ = fmt.Printf
var _ = strings.Trim

type Flag struct {
    names []string
    has_arg bool
}

type ShortFlag Flag
type LongFlag Flag

type PosArg string

type MatchResult map[string]string

type Opt interface {
    Match(string, MatchResult) (MatchResult, error)
}

func ParseOpt(line string) (res Opt, err error) {
    match := ReOptDesc.FindStringSubmatch(line)

    return nil, nil
}

type Group struct {
    items []Opt
    optional bool
}

type ExclusiveGroup Group

var ReSections = re.MustCompile(
    `(?is)usage: *(.+?)` +
    `(?:\n[ \t]*\n(?: *options: *\n)?(.+?))?` +
    `(?:\n[ \t]*\n.*)?$`)
var RePosArg = re.MustCompile(
    `(<[^>]+>|[[:upper:]]+)`)
var ReFlagName = re.MustCompile(
    `(--[^ =<]+|-[^ =<])`)
var ReFlag = re.MustCompile(
    ReFlagName.String() + `(?:[= ]?` + RePosArg.String() + `)?`)
var ReOptDescDelim = re.MustCompile(
    `(?:, *| +)`)
var ReOptDesc = re.MustCompile(
    `((?:` + ReFlag.String() + ReOptDescDelim.String() + `?)+)` +
    ` {2}` +
    `(.+)`)

func Docopt(doc string) (res MatchResult, err error) {
    sections := ReSections.FindStringSubmatch(doc)
    if len(sections) < 3 {
        panic(`"Usage:" (case-insensetive) not found in help text`)
    }

    parseOptions(sections[2])

    return nil, nil
}

func parseOptions(opts string) {
    if len(opts) == 0 {
        return
    }

    for _, line := range strings.Split(opts, "\n") {
        line = strings.Trim(line, " \t")
        ParseOpt(line)
    }
}
