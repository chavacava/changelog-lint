// Configuration for changelog-lint
package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/linting/rule"
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

// RulesConfig defines the config for all rules.
type RulesConfig = map[string]RuleConfig

type Config struct {
	RulesConfs RulesConfig `toml:"rule"`
}

func (c Config) enabledRules() []linting.Rule {
	availableRules := make(map[string]linting.Rule, len(allRules))
	for _, r := range allRules {
		availableRules[r.Name()] = r
	}

	enabledRules := []linting.Rule{}
	for ruleName, r := range availableRules {
		if c.RulesConfs[ruleName].Disabled {
			continue
		}

		enabledRules = append(enabledRules, r)
	}

	return enabledRules
}

func defaultConfig() *linting.Config {
	config := &linting.Config{
		RuleArgs: map[linting.Rule]any{},
	}

	for _, r := range allRules {
		config.RuleArgs[r] = nil
	}

	return config
}

func GetConfig(configFile string) (*linting.Config, error) {
	if configFile == "" {
		return defaultConfig(), nil
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading the config file %s: %v", configFile, err)
	}
	config := &Config{}
	_, err = toml.Decode(string(file), config)
	if err != nil {
		return nil, fmt.Errorf("error parsing the config file %s: %v", configFile, err)
	}

	result := &linting.Config{RuleArgs: map[linting.Rule]any{}}
	for _, r := range config.enabledRules() {
		result.RuleArgs[r] = config.RulesConfs[r.Name()].Arguments
	}

	return result, nil
}
