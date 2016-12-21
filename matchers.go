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

	MatcherOptionNameSet = NewMatcher(
		`[^ =<|\])]`,
	)

	MatcherOptionName = NewMatcher(
		`(--%[1]s+|-%[1]s)`,
		MatcherOptionNameSet,
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

	MatcherDescriptionParagraph = NewMatcher(
		`\W|^$`,
	)

	MatcherDescriptionDefault = NewMatcher(
		`(?s)(?:.*)\[default: ([^\]]+)]`,
	)

	MatcherTokenSeparator = NewMatcher(
		`\s+`,
	)

	MatcherEndOfLine = NewMatcher(
		`$`,
	)

	MatcherTokenBranch = NewMatcher(
		`\|`,
	)

	MatcherTokenRequiredGroupStart = NewMatcher(
		`\(`,
	)

	MatcherTokenOptionalGroupStart = NewMatcher(
		`\[`,
	)

	MatcherTokenRequiredGroupEnd = NewMatcher(
		`\)`,
	)

	MatcherTokenOptionalGroupEnd = NewMatcher(
		`\]`,
	)

	MatcherTokenRepeat = NewMatcher(
		`\.\.\.`,
	)

	MatcherTokenWord = NewMatcher(
		`([^ \t([\])<>|.-]%[1]s*)`,
		MatcherOptionNameSet,
	)
)
