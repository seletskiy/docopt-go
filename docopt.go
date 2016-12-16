package docopt

var (
	MatcherSections = NewMatcher(
		`(?is)` +
			`usage: *\n(.+?)` +
			`(?:\n[ \t]*\n(?: *[^:\n]+: *)?\n(.+?))?` +
			`(?:\n[ \t]*\n.*)?$`,
	)

	MatcherArgument = NewMatcher(
		`(<[^>]+>|[[:upper:]]+)`,
	)

	MatcherOptionName = NewMatcher(
		`(--[^ =<]+|-[^ =<])`,
	)

	MatcherOption = NewMatcher(
		`%[1]s(?:[= ]?%[2]s)?`,
		MatcherOptionName,
		MatcherArgument,
	)

	MatcherOptionSeparator = NewMatcher(
		`(?:, *| +)`,
	)

	MatcherDescriptionSeparator = NewMatcher(
		` {2,}`,
	)

	MatcherIndenting = NewMatcher(
		`([ \t]*)`,
	)
)
