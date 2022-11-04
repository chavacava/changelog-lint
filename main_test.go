package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		args []string
		want int
	}{
		{
			args: []string{"changelog-lint", "-version"},
			want: codeOK,
		},
		{
			args: []string{"changelog-lint"},
			want: codeOK,
		},
		{
			args: []string{"changelog-lint", "unknownfile.md"},
			want: codeRequestError,
		},
		{
			args: []string{"changelog-lint", "-config", "unknownConfig.toml"},
			want: codeRequestError,
		},
		{
			args: []string{"changelog-lint", "-config", "testdata/parser-conf-bad-pattern-entry.toml"},
			want: codeRequestError,
		},
		{
			args: []string{"changelog-lint", "-config", "testdata/malformed.toml"},
			want: codeRequestError,
		},
		{
			args: []string{"changelog-lint", "-release", "0.0.0"},
			want: codeLintError,
		},
		{
			args: []string{"changelog-lint", "./testdata/parser-error.md"},
			want: codeSyntaxError,
		},
		{
			args: []string{"changelog-lint", "./testdata/keepachangelog.md"},
			want: codeOK,
		},
	}
	for _, tc := range testCases {
		got := run(tc.args)
		if got != tc.want {
			t.Errorf("expected %d for %v, got %d", tc.want, tc.args, got)
		}
	}

}
