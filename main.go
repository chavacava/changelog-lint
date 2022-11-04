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
	os.Exit(run(os.Args))
}

func run(args []string) int {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagVersion := flags.Bool("version", false, "get changelog-lint version")
	flagConfig := flags.String("config", "", "set linter configuration")
	flagReleaseMode := flags.String("release", "", "enables release-related checks (the given string must be the release version, e.g. 1.2.3)")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Println(err)
		return codeRequestError
	}

	if *flagVersion {
		println(versionInfo())
		return codeOK
	}

	freeArgs := flags.Args()
	inputFilename := "CHANGELOG.md"
	if len(freeArgs) > 0 {
		inputFilename = freeArgs[0]
	}
	input, err := os.Open(inputFilename)
	if err != nil {
		fmt.Println(err)
		return codeRequestError
	}
	defer input.Close()

	mainConfig, err := config.LoadConfig(*flagConfig)
	if err != nil {
		fmt.Println(err)
		return codeRequestError
	}

	parserConf, err := mainConfig.ParserConfig()
	if err != nil {
		fmt.Println(err)
		return codeRequestError
	}

	p := parser.Default{}
	changes, err := p.Parse(input, parserConf)
	if err != nil {
		fmt.Println(err)
		return codeSyntaxError
	}

	linter := linting.Linter{}
	failures := make(chan linting.Failure)
	lintingConfig := mainConfig.LintingConfig()
	if *flagReleaseMode != "" {
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
	return exitCode
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
