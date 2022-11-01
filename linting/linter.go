// Linting package aggregates the linting machinery
package linting

import (
	"github.com/chavacava/changelog-lint/model"
)

type RuleArgs = any

// Rule general abstraction
type Rule interface {
	Apply(model.Changelog, chan Failure, RuleArgs)
	Name() string
}

// Linter provides changelog linting method
type Linter struct{}

// Lint a changelog
func (Linter) Lint(changes model.Changelog, config *Config, failures chan Failure) {
	for rule, rConf := range config.RuleArgs {
		rule.Apply(changes, failures, rConf)
	}
	close(failures)
}
