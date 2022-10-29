package rule

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/chavacava/changelogger/linting"
	"github.com/chavacava/changelogger/model"
	"github.com/chavacava/changelogger/parser"
)

type testBundle struct {
	rule      linting.Rule
	testCases map[string][]string // map[filename][]expectedFailureMessages
}

var bundles []testBundle = []testBundle{
	{
		SubsectionEmpty{},
		map[string][]string{
			"subection-empty.md": {
				`empty subsection "### Infrastructure" in version ## 1.9.0 - 1711-04-20`,
				`empty subsection "### Added" in version ## 1.9.0 - 1711-04-20`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
	{
		VersionEmpty{},
		map[string][]string{
			"version-empty.md": {
				`empty version "## 1.9.0"`,
				`empty version "## 1.6.0"`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
	{
		SubsectionNamming{},
		map[string][]string{
			"subsection-namming.md": {
				`unknown subsection "### Infrastructures" in version ## 1.9.0 - 1111-04-20`,
				`unknown subsection "### Addeda" in version ## 1.9.0 - 1111-04-20`,
			},
			"ok.md": { /* no error expected */ },
		},
	},
}

func TestRules(t *testing.T) {
	for _, bundle := range bundles {
		for file, wantFailureMessages := range bundle.testCases {
			changes, err := parseChangelog(file)
			if err != nil {
				t.Fatal(err.Error())
			}

			if err := ruleTester(bundle.rule, *changes, wantFailureMessages); err != nil {
				t.Fatal(err.Error())
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
	changes, err := p.Parse(input, map[string]string{})
	if err != nil {
		return nil, err
	}

	return changes, nil
}

func ruleTester(r linting.Rule, changes model.Changelog, wantFailureMessages []string) error {
	failures := make(chan linting.Failure)

	go func() {
		r.Apply(changes, failures, nil)
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
