package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestLoadConfigRulesPart(t *testing.T) {
	// Default conf retrieval
	got, err := LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range allRules {
		rc, ok := got.Rules[r.Name()]
		if !ok {
			t.Fatalf("rule %s not present in the conf %+v", r.Name(), got.Rules)
		}

		if len(rc.Arguments) != 0 || rc.Disabled {
			t.Fatalf("expected no conf for rule %s, got %+v", r.Name(), rc)
		}
	}

	// Error finding the conf file
	_, err = LoadConfig("unknown-conf-file")
	if err == nil {
		t.Fatal("expected 'conf file not found' error but got nothing")
	}

	// Error parsing the conf file
	_, err = LoadConfig("./testdata/malformed.toml")
	if err == nil {
		t.Fatal("expected conf file parsing error but got no error")
	}

	// Load real conf file
	got, err = LoadConfig("./testdata/conf.toml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range allRules {
		rc, ok := got.Rules[r.Name()]

		// Check disabled rule
		switch r.Name() {
		case "version-empty":
			if ok {
				t.Fatal("rule \"version-empty\" should be disabled")
			}
		default:
			if !ok {
				t.Fatalf("rule %s is not present in the conf %+v", r.Name(), got.LintingConfig().RuleArgs)
			}
		}

		if !ok {
			continue
		}

		args := fmt.Sprintf("%v", rc.Arguments)
		// Check rule args
		switch r.Name() {
		case "version-order":
			want := "[1 2 3 4]"
			if args != want {
				t.Fatalf("expected conf for rule %s to be %s, got %s", r.Name(), want, args)
			}
		case "version-repetition":
			want := "[a string another one]"
			if args != want {
				t.Fatalf("expected conf for rule %s to be %s, got %s", r.Name(), want, args)
			}
		default:
			if args != "[]" {
				t.Fatalf("expected empty conf for rule %s, got %s", r.Name(), args)
			}
		}
	}
}

func TestLoadConfigParserPart(t *testing.T) {
	got, err := LoadConfig("./testdata/parser-conf.toml")
	if err != nil {
		t.Fatalf("unexpected conf file parsing error: %v", err)
	}

	wantPattern := "title pattern"
	gotPattern := got.Parser.Patterns.Title
	if gotPattern != wantPattern {
		t.Fatalf("expected pattern to be %s, got %s", wantPattern, gotPattern)
	}

	wantPattern = "version pattern"
	gotPattern = got.Parser.Patterns.Version
	if gotPattern != wantPattern {
		t.Fatalf("expected pattern to be %s, got %s", wantPattern, gotPattern)
	}

	wantPattern = "subsection pattern"
	gotPattern = got.Parser.Patterns.Subsection
	if gotPattern != wantPattern {
		t.Fatalf("expected pattern to be %s, got %s", wantPattern, gotPattern)
	}

	wantPattern = "entry pattern"
	gotPattern = got.Parser.Patterns.Entry
	if gotPattern != wantPattern {
		t.Fatalf("expected pattern to be %s, got %s", wantPattern, gotPattern)
	}
}

func TestParserConfigErrors(t *testing.T) {
	testCases := []struct {
		file string
		want string
	}{
		{
			file: "./testdata/parser-conf-bad-pattern-title.toml",
			want: "title's pattern",
		},
		{
			file: "./testdata/parser-conf-bad-pattern-version.toml",
			want: "version's pattern",
		},
		{
			file: "./testdata/parser-conf-bad-pattern-subsection.toml",
			want: "subsection's pattern",
		},
		{
			file: "./testdata/parser-conf-bad-pattern-entry.toml",
			want: "entry's pattern",
		},
	}

	for _, tc := range testCases {
		config, err := LoadConfig(tc.file)
		if err != nil {
			t.Fatalf("unexpected config error: %v", err)
		}

		_, err = config.ParserConfig()
		if err == nil {
			t.Fatal("expected config error, got nothing")
		}

		got := err.Error()
		if !strings.HasPrefix(got, tc.want) {
			t.Fatalf("expected config error to be %q... got %q", tc.want, got)
		}
	}
}

func TestParserConfig(t *testing.T) {
	config, err := LoadConfig("./testdata/parser-conf.toml")
	if err != nil {
		t.Fatalf("unexpected config error: %v", err)
	}

	got, err := config.ParserConfig()
	if err != nil {
		t.Fatalf("unexpected config error: %v ", err)
	}

	wantPattern := "title pattern"
	gotPattern := got.TitlePattern.String()
	if gotPattern != wantPattern {
		t.Fatalf("expected pattern to be %s, got %s", wantPattern, gotPattern)
	}

	wantPattern = "version pattern"
	gotPattern = got.VersionPattern.String()
	if gotPattern != wantPattern {
		t.Fatalf("expected pattern to be %s, got %s", wantPattern, gotPattern)
	}

	wantPattern = "subsection pattern"
	gotPattern = got.SubsectionPattern.String()
	if gotPattern != wantPattern {
		t.Fatalf("expected pattern to be %s, got %s", wantPattern, gotPattern)
	}

	wantPattern = "entry pattern"
	gotPattern = got.EntryPattern.String()
	if gotPattern != wantPattern {
		t.Fatalf("expected pattern to be %s, got %s", wantPattern, gotPattern)
	}
}
