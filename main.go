package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/linting/rule"
	"github.com/chavacava/changelog-lint/parser"
)

func main() {
	flag.Parse()
	args := flag.Args()

	inputFilename := "CHANGELOG.md"
	if len(args) > 0 {
		inputFilename = args[0]
	}
	input, err := os.Open(inputFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := parser.Default{}
	changes, err := p.Parse(input, map[string]string{})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
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
		exitCode = 3
	}

	os.Exit(exitCode)
}
