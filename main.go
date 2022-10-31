package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/linting/rule"
	"github.com/chavacava/changelog-lint/parser"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

const (
	codeOK = iota
	codeRequestError
	codeSyntaxError
	codeLintError
)

func main() {
	flagVersion := flag.Bool("version", false, "get changelog-lint version")
	flag.Parse()

	if *flagVersion {
		println(versionInfo())
		os.Exit(codeOK)
	}
	args := flag.Args()
	inputFilename := "CHANGELOG.md"
	if len(args) > 0 {
		inputFilename = args[0]
	}
	input, err := os.Open(inputFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(codeRequestError)
	}
	defer input.Close()

	p := parser.Default{}
	changes, err := p.Parse(input, map[string]string{})
	if err != nil {
		fmt.Println(err)
		os.Exit(codeSyntaxError)
	}

	linter := linting.Linter{}
	failures := make(chan linting.Failure)
	r1 := rule.SubsectionNamming{}
	r2 := rule.SubsectionOrder{}
	r3 := rule.VersionOrder{}
	r4 := rule.VersionRepetition{}
	r5 := rule.VersionEmpty{}

	config := linting.Config{
		Rules:     []linting.Rule{r1, r2, r3, r4, r5},
		RuleConfs: map[string][]any{},
	}

	go linter.Lint(*changes, &config, failures)

	exitCode := codeOK

	for failure := range failures {
		lineInfo := ""
		if failure.Position > 0 {
			lineInfo = fmt.Sprintf("(line %d)", failure.Position)
		}
		fmt.Printf("%s: %s %s\n", failure.RuleName, failure.Message, lineInfo)
		exitCode = codeLintError
	}

	os.Exit(exitCode)
}

func versionInfo() string {
	var buildInfo string
	if date != "unknown" && builtBy != "unknown" {
		buildInfo = fmt.Sprintf("Built\t\t%s by %s\n", date, builtBy)
	}

	if commit != "none" {
		buildInfo = fmt.Sprintf("Commit:\t\t%s\n%s", commit, buildInfo)
	}

	if version == "dev" {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			version = bi.Main.Version
			if strings.HasPrefix(version, "v") {
				version = bi.Main.Version[1:]
			}
			if len(buildInfo) == 0 {
				return fmt.Sprintf("version %s\n", version)
			}
		}
	}

	return fmt.Sprintf("Version:\t%s\n%s", version, buildInfo)
}
