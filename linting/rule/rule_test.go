package rule

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
	"github.com/chavacava/changelog-lint/parser"
)

type testBundle struct {
	rule      linting.Rule
	args      any
	testCases map[string][]string // map[filename][]expectedFailureMessages
}

var bundles []testBundle = []testBundle{
	{
		SubsectionEmpty{},
		nil,
		map[string][]string{
			"subection-empty.md": {
				`empty subsection "Infrastructure" in version 1.9.0`,
				`empty subsection "Added" in version 1.9.0`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
	{
		VersionEmpty{},
		nil,
		map[string][]string{
			"version-empty.md": {
				`empty version "1.9.0"`,
				`empty version "1.6.0"`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
	{
		SubsectionNaming{},
		nil,
		map[string][]string{
			"subsection-naming.md": {
				`unknown subsection "Infrastructures" in version 1.9.0`,
				`unknown subsection "Addeda" in version 1.9.0`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
	{
		SubsectionOrder{},
		nil,
		map[string][]string{
			"subsection-order.md": {
				`subsection "Added" is not sorted alphabetically in version 1.9.0`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
	{
		SubsectionRepetition{},
		nil,
		map[string][]string{
			"subsection-repetition.md": {
				`duplicated subsection "Added"`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
	{
		VersionRepetition{},
		nil,
		map[string][]string{
			"version-repetition.md": {
				`duplicated version 1.8.0`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
	{
		VersionOrder{},
		nil,
		map[string][]string{
			"version-order.md": {
				`version Unreleased must be at the top of the version list`,
				`version 1.10.0 is not well sorted`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
	{
		Release{},
		nil,
		map[string][]string{
			"ok.md": {
				"expected release version argument to be a string, got <nil> instead",
			},
		},
	},
	{
		Release{},
		"1.0.0",
		map[string][]string{
			"ok.md": {
				"expected release version to be 1.0.0, got Unreleased instead",
				"version Unreleased forbidden in release mode",
			},
		},
	},
}

func TestRules(t *testing.T) {
	for _, bundle := range bundles {
		for file, wantFailureMessages := range bundle.testCases {
			changes, err := parseChangelog(file)
			if err != nil {
				t.Fatalf("%s/%s:%v", bundle.rule.Name(), file, err)
			}

			if err := ruleTester(bundle.rule, bundle.args, *changes, wantFailureMessages); err != nil {
				t.Fatalf("%s/%s:%v", bundle.rule.Name(), file, err)
			}
		}
	}
}

func parseChangelog(filename string) (*model.Changelog, error) {
	input, err := os.Open(filepath.Join("testdata", filename))
	if err != nil {
		return nil, err
	}

	p := parser.Default{}
	changes, err := p.Parse(input, parserConf())
	if err != nil {
		return nil, err
	}

	return changes, nil
}

func ruleTester(r linting.Rule, args any, changes model.Changelog, wantFailureMessages []string) error {
	failures := make(chan linting.Failure)

	go func() {
		r.Apply(changes, failures, args)
		close(failures)
	}()

	i := 0
	for failure := range failures {
		if i >= len(wantFailureMessages) {
			return fmt.Errorf("unexpected failure:\n\t%v", failure)
		}
		wantFailure := wantFailureMessages[i]
		got := failure.Message
		if got != wantFailure {
			return fmt.Errorf("expected failure:\n\t%q\ngot:\n\t%q", wantFailure, got)
		}
		i++
	}
	if i < len(wantFailureMessages) {
		msg := "unsatisfied expected failure(s):"
		for ; i < len(wantFailureMessages); i++ {
			msg += "\n" + wantFailureMessages[i]
		}
		return errors.New(msg)
	}

	return nil
}

func parserConf() *parser.Config {
	return &parser.Config{
		TitlePattern:      regexp.MustCompile(`.+`),
		VersionPattern:    regexp.MustCompile(`^## \[?(\d+\.\d+.\d+|Unreleased)\]?( .*)*$`),
		SubsectionPattern: regexp.MustCompile(`^### ([A-Z]+[a-z]+)[ ]*$`),
		EntryPattern:      regexp.MustCompile(`^[*-] .+$`),
	}
}
