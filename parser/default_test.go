package parser_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/chavacava/changelog-lint/parser"
)

func TestDefaultParser(t *testing.T) {
	testsCases := []struct {
		data string
		err  string
	}{
		{
			data: "CHANGELOG_OK_0.md",
			err:  "",
		},
		{
			data: "CHANGELOG_OK_1.md",
			err:  "",
		},
		{
			data: "CHANGELOG_OK_2.md",
			err:  "",
		},
		{
			data: "CHANGELOG_ERR_1.md",
			err:  "unexpected end of file",
		},
		{
			data: "CHANGELOG_ERR_3.md",
			err:  "unexpected end of file",
		},
		{
			data: "CHANGELOG_ERR_4.md",
			err: `unexpected line:some nice feature
expecting subsection, version or change description`,
		},
		{
			data: "CHANGELOG_ERR_6.md",
			err:  "the version\n\t## 19.0\n" + `does not match ^## \[?(\d+\.\d+.\d+|Unreleased)\]?( .*)*$ (line 7)`,
		},
		{
			data: "CHANGELOG_ERR_7.md",
			err:  "version \"19.00.0\" contains a bad formated subsection: the subsection\n\t### 12 Adds\n" + `does not match ^### ([A-Z]+[a-z]+)[ ]*$ (line 9)`,
		},
	}

	for _, tc := range testsCases {
		p := parser.Default{}
		dataPath := "testdata/" + tc.data
		file, err := os.Open(dataPath)
		if err != nil {
			t.Fatalf("error reading test data: %v", err)
		}

		_, err = p.Parse(file, parserConf())

		switch {
		case err != nil && tc.err == "":
			// unexpected error
			t.Fatalf("unexpected error parsing %q:\n%v", dataPath, err)
		case err != nil && tc.err != "":
			// check error msg
			if !strings.HasPrefix(err.Error(), tc.err) {
				t.Fatalf("expected error parsing %q:\n%v\ngot:\n%s", dataPath, tc.err, err)
			}
		case err == nil && tc.err != "":
			// should fail
			t.Fatalf("expected error\n%q\nparsing %q:\n", tc.err, dataPath)
		default:
			// OK
		}
	}
}

func parserConf() *parser.Config {
	return &parser.Config{
		TitlePattern:      regexp.MustCompile(`.+`),
		VersionPattern:    regexp.MustCompile(`^## \[?(\d+\.\d+.\d+|Unreleased)\]?( .*)*$`),
		SubsectionPattern: regexp.MustCompile(`^### ([A-Z]+[a-z]+)[ ]*$`),
		EntryPattern:      regexp.MustCompile(`^[*-] .+$`),
	}
}
