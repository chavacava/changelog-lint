package model

import (
	"strings"
)

// Highly inspired from github.com/parkr/changelog

// Changelog represents a changelog in its entirety, containing all the
// versions that are tracked in the changelog. For supported formats, see
// the documentation for Version.
type Changelog struct {
	Header   []string
	Versions []*Version
}

// A Markdown string representation of the Changelog.
func (c *Changelog) String() string {
	//c.sortVersions() TODO
	versionStrs := make([]string, len(c.Versions))
	for i, version := range c.Versions {
		versionStrs[i] = version.String()
	}
	return strings.Join(versionStrs, "\n\n") + "\n"
}

// Version contains the data for the changes for a given version. It can
// have both direct history and subsections.
// Acceptable formats:
//
//	## 2.4.1
//	## 2.4.1 / 2015-04-23
//
// The version currently cannot be prefixed with a `v`, but a date is
// optional.
type Version struct {
	Version     string
	Subsections []*Subsection
}

// String returns the markdown representation for the version.
func (v *Version) String() string {
	lines := []string{}
	if v.Version != "" {
		if len(lines) == 0 {
			lines = append(lines, "")
		}
		lines[0] += "## " + v.Version
	}
	if len(v.Subsections) > 0 {
		subsectionsStrs := make([]string, len(v.Subsections))
		for i, subsection := range v.Subsections {
			subsectionsStrs[i] = subsection.String()
		}
		lines = append(lines, strings.Join(subsectionsStrs, "\n\n"))
	}
	return strings.Join(lines, "\n\n")
}

// Subsection contains the data for a given subsection.
// Acceptable format:
//
//	### Subsection Name Here
//
// Common subsections are "Major Enhancements," and "Bug Fixes."
type Subsection struct {
	Name    string
	History []*ChangeLine
}

// String returns the markdown representation of the subsection.
func (s *Subsection) String() string {
	if len(s.History) > 0 {
		historyStrs := make([]string, len(s.History))
		for i, history := range s.History {
			historyStrs[i] = history.String()
		}
		return "### " + s.Name + "\n\n" + strings.Join(historyStrs, "\n")
	}
	return ""
}

// ChangeLine contains the data for a single change.
// Acceptable formats:
//
//   - This is a change (#1234)
//   - This is another change. (@parkr)
//   - This is a change w/o a reference.
//
// The references must be encased in parentheses, and only one reference is
// currently supported.
type ChangeLine struct {
	// What the change entails.
	Summary string
}

// String returns the markdown representation of the ChangeLine.
// E.g. "  * Added documentation. (#123)"
func (l *ChangeLine) String() string {
	str := "  * " + l.Summary
	return str
}

// NewChangelog creates a pristine Changelog.
func NewChangelog() *Changelog {
	return &Changelog{Versions: []*Version{}}
}
