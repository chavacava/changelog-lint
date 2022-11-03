// Configuration for changelog-lint
package config

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/linting/rule"
	"github.com/chavacava/changelog-lint/parser"
)

var allRules = []linting.Rule{
	rule.SubsectionEmpty{},
	rule.SubsectionNamming{},
	rule.SubsectionOrder{},
	rule.SubsectionRepetition{},
	rule.VersionEmpty{},
	rule.VersionOrder{},
	rule.VersionRepetition{},
}

// Arguments is type used for the arguments of a rule.
type Arguments = []any

// RuleConfig is type used for the rule configuration.
type RuleConfig struct {
	Arguments Arguments
	Disabled  bool
}

type ParserPatterns struct {
	Title      string
	Version    string
	Subsection string
	Entry      string
}

// RuleConfig is type used for the rule configuration.
type ParserConfig struct {
	Patterns ParserPatterns
}

// RulesConfig defines the config for all rules.
type RulesConfig = map[string]RuleConfig

type Config struct {
	Rules  RulesConfig  `toml:"rule"`
	Parser ParserConfig `toml:"parser"`
}

func (c Config) enabledRules() []linting.Rule {
	availableRules := make(map[string]linting.Rule, len(allRules))
	for _, r := range allRules {
		availableRules[r.Name()] = r
	}

	enabledRules := []linting.Rule{}
	for ruleName, r := range availableRules {
		if c.Rules[ruleName].Disabled {
			continue
		}

		enabledRules = append(enabledRules, r)
	}

	return enabledRules
}

func defaultRulesConfig() RulesConfig {
	config := make(map[string]RuleConfig, len(allRules))

	for _, r := range allRules {
		config[r.Name()] = RuleConfig{}
	}

	return config
}

func defaultLintingConfig() *linting.Config {
	config := &linting.Config{
		RuleArgs: map[linting.Rule]any{},
	}

	for _, r := range allRules {
		config.RuleArgs[r] = nil
	}

	return config
}

const defaultPatternTitle = `.+`
const defaultPatternVersion = `^## \[?(\d+\.\d+.\d+|Unreleased)\]?( .*)*$`
const defaultPatternSubsection = `^### ([A-Z]+[a-z]+)[ ]*$`
const defaultPatternEntry = `^[*-] .+$`

func defaultParseConf() ParserConfig {
	return ParserConfig{
		Patterns: ParserPatterns{
			Title:      defaultPatternTitle,
			Version:    defaultPatternVersion,
			Subsection: defaultPatternSubsection,
			Entry:      defaultPatternEntry,
		},
	}
}

func defaultConf() *Config {
	return &Config{
		Rules:  defaultRulesConfig(),
		Parser: defaultParseConf(),
	}
}

func LoadConfig(configFile string) (*Config, error) {
	defaultConf := defaultConf()
	if configFile == "" {
		return defaultConf, nil
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading the config file %s: %v", configFile, err)
	}
	loadedConf := &Config{}
	_, err = toml.Decode(string(file), loadedConf)
	if err != nil {
		return nil, fmt.Errorf("error parsing the config file %s: %v", configFile, err)
	}

	for k, v := range loadedConf.Rules {
		if v.Disabled {
			delete(defaultConf.Rules, k)
			continue
		}
		defaultConf.Rules[k] = v
	}

	if loadedConf.Parser.Patterns.Title != "" {
		defaultConf.Parser.Patterns.Title = loadedConf.Parser.Patterns.Title
	}

	if loadedConf.Parser.Patterns.Version != "" {
		defaultConf.Parser.Patterns.Version = loadedConf.Parser.Patterns.Version
	}

	if loadedConf.Parser.Patterns.Subsection != "" {
		defaultConf.Parser.Patterns.Subsection = loadedConf.Parser.Patterns.Subsection
	}

	if loadedConf.Parser.Patterns.Entry != "" {
		defaultConf.Parser.Patterns.Entry = loadedConf.Parser.Patterns.Entry
	}

	return defaultConf, nil
}

func (c Config) LintingConfig() *linting.Config {
	if c.Rules == nil {
		return defaultLintingConfig()
	}

	result := &linting.Config{RuleArgs: map[linting.Rule]any{}}
	for _, r := range c.enabledRules() {
		result.RuleArgs[r] = c.Rules[r.Name()].Arguments
	}

	return result
}

func (c Config) ParserConfig() (*parser.Config, error) {
	reTitle, err := regexp.Compile(c.Parser.Patterns.Title)
	if err != nil {
		return nil, fmt.Errorf("title's pattern %s in configuration does not compile: %v", c.Parser.Patterns.Title, err)
	}
	reVersion, err := regexp.Compile(c.Parser.Patterns.Version)
	if err != nil {
		return nil, fmt.Errorf("version's pattern %s in configuration does not compile: %v", c.Parser.Patterns.Version, err)
	}
	reSubsection, err := regexp.Compile(c.Parser.Patterns.Subsection)
	if err != nil {
		return nil, fmt.Errorf("subsection's pattern %s in configuration does not compile: %v", c.Parser.Patterns.Subsection, err)
	}
	reEntry, err := regexp.Compile(c.Parser.Patterns.Entry)
	if err != nil {
		return nil, fmt.Errorf("entry's pattern %s in configuration does not compile: %v", c.Parser.Patterns.Entry, err)
	}
	result := &parser.Config{
		TitlePattern:      reTitle,
		VersionPattern:    reVersion,
		SubsectionPattern: reSubsection,
		EntryPattern:      reEntry,
	}
	return result, nil
}
