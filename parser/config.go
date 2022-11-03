package parser

import "regexp"

// Config for parsers
type Config struct {
	TitlePattern      *regexp.Regexp
	VersionPattern    *regexp.Regexp
	SubsectionPattern *regexp.Regexp
	EntryPattern      *regexp.Regexp
}
