package linting

import (
	"github.com/chavacava/changelogger/model"
)

type Arguments = []any

type Rule interface {
	Apply(model.Changelog, chan Failure, Arguments)
	Name() string
}
type Linter struct{}

func (Linter) Lint(changes model.Changelog, conf *Config, failures chan Failure) {
	for _, rule := range conf.Rules {
		ruleConfig := conf.RuleConfs[rule.Name()]
		rule.Apply(changes, failures, ruleConfig)
	}
	close(failures)
}
