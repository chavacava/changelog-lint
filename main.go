package main

import (
	"os"

	"github.com/chavacava/changelogger/linting"
	"github.com/chavacava/changelogger/linting/rule"
	"github.com/chavacava/changelogger/parser"
)

func main() {
	input, err := os.Open(("CHANGELOG.md"))
	if err != nil {
		panic(err)
	}

	p := parser.Default{}
	changes, err := p.Parse(input, map[string]string{})
	if err != nil {
		panic(err)
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

	exitCode := 0

	for failure := range failures {
		println(failure.Originator + ": " + failure.Message)
		exitCode = 1
	}

	os.Exit(exitCode)
}
