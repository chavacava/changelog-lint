package config

import (
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	// Default conf retrieval
	got, err := GetConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range allRules {
		rc, ok := got.RuleArgs[r]
		if !ok {
			t.Fatalf("rule %s not present in the conf %+v", r.Name(), got.RuleArgs)
		}
		if rc != nil {
			t.Fatalf("expected nil conf for rule %s, got %+v", r.Name(), rc)
		}
	}

	// Error finding the conf file
	_, err = GetConfig("unknown-conf-file")
	if err == nil {
		t.Fatal("expected 'conf file not found' error but got nothing")
	}

	// Error parsing the conf file
	got, err = GetConfig("./testdata/conf.toml")
	if err != nil {
		t.Fatalf("unexpected conf file parsing error: %v", err)
	}

	// Load real conf file
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range allRules {
		rc, ok := got.RuleArgs[r]

		// Check disabled rule
		switch r.Name() {
		case "version-empty":
			if ok {
				t.Fatal("rule \"version-empty\" should be disabled")
			}
		default:
			if !ok {
				t.Fatalf("rule %s is not present in the conf %+v", r.Name(), got.RuleArgs)
			}
		}

		if !ok {
			continue
		}

		conf := rc.(RuleConfig)
		args := fmt.Sprintf("%v", conf.Arguments)
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
