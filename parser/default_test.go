package parser_test

import (
	"os"
	"testing"

	"github.com/chavacava/changelogger/parser"
)

func TestDefaultParser(t *testing.T) {
	testsCases := []struct {
		data string
		err  string
	}{
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
	}

	for _, tc := range testsCases {
		p := parser.Default{}
		dataPath := "testdata/" + tc.data
		file, err := os.Open(dataPath)
		if err != nil {
			t.Fatalf("error reading test data: %v", err)
		}

		_, err = p.Parse(file, map[string]string{})

		switch {
		case err != nil && tc.err == "":
			// unexpected error
			t.Fatalf("unexpected error parsing %q:\n%v", dataPath, err)
		case err != nil && tc.err != "":
			// check error msg
			if err.Error() != tc.err {
				t.Fatalf("expected error parsing %q:\n%v\ngot:\n%s", dataPath, tc.err, err)
			}
		case err == nil && tc.err != "":
			// should faild
			t.Fatalf("expected error\n%q\nparsing %q:\n", tc.err, dataPath)
		default:
			// OK
		}
	}
}
