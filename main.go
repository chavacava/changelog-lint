package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/chavacava/changelog-lint/config"
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
	flagConfig := flag.String("config", "", "set linter configuration")
	flagReleaseMode := flag.String("release", "", "enables release-related checks (the given string must be the release version, e.g. 1.2.3)")
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

	mainConfig, err := config.LoadConfig(*flagConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(codeRequestError)
	}

	parserConf, err := mainConfig.ParserConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(codeRequestError)
	}

	p := parser.Default{}
	changes, err := p.Parse(input, parserConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(codeSyntaxError)
	}

	linter := linting.Linter{}
	failures := make(chan linting.Failure)
	lintingConfig := mainConfig.LintingConfig()
	if *flagReleaseMode != "" { // Add release rule to the linting conf
		releaseRule := &rule.Release{}
		lintingConfig.RuleArgs[releaseRule] = *flagReleaseMode
	}
	go linter.Lint(*changes, lintingConfig, failures)

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
