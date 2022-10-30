package model

// Highly inspired from github.com/parkr/changelog

// Changelog represents a changelog in its entirety, containing all the
// versions that are tracked in the changelog. For supported formats, see
// the documentation for Version.
type Changelog struct {
	Header   []string
	Versions []*Version
}

// NewChangelog creates a pristine Changelog.
func NewChangelog() *Changelog {
	return &Changelog{Versions: []*Version{}}
}

// Version contains the data for the changes for a given version. It can
// have both direct history and subsections.
type Version struct {
	Version     string
	Subsections []*Subsection
	SourceLine  string
	Position    int // Line number in the changelog
}

// Subsection contains the data for a given subsection.
type Subsection struct {
	Name       string
	History    []*Entry
	SourceLine string
	Position   int // Line number in the changelog
}

// Entry contains the data for a single change.
type Entry struct {
	// What the change entails.
	Summary  string
	Position int // Line number in the changelog
}
