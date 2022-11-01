// Linting package aggregates the linting machinery
package linting

import (
	"github.com/chavacava/changelog-lint/model"
)

type Arguments = []any

// Rule general abstraction
type Rule interface {
	Apply(model.Changelog, chan Failure, Arguments)
	Name() string
}

// Linter provides changelog linting method
type Linter struct{}

// Lint a changelog
func (Linter) Lint(changes model.Changelog, conf *Config, failures chan Failure) {
	for _, rule := range conf.Rules {
		ruleConfig := conf.RuleConfs[rule.Name()]
		rule.Apply(changes, failures, ruleConfig)
	}
	close(failures)
}
