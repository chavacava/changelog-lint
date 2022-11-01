package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type SubsectionNamming struct{}

func (r SubsectionNamming) Apply(changes model.Changelog, failures chan linting.Failure, args linting.RuleConf) {
	allowedSubsections, err := r.allowedSubsections(args)
	if err != nil {
		msg := fmt.Sprintf("bad rule configuration for %q: %v", r.Name(), err)
		failures <- linting.Failure{RuleName: r.Name(), Message: msg}
		return
	}

	for _, version := range changes.Versions {
		for _, subsection := range version.Subsections {
			_, ok := allowedSubsections[subsection.Name]
			if ok {
				continue
			}
			msg := fmt.Sprintf("unknown subsection %q in version %v", subsection.Name, version.Version)
			failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: subsection.Position}
		}
	}
}

func (SubsectionNamming) Name() string {
	return "subsection-namming"
}

func (r SubsectionNamming) allowedSubsections(args linting.RuleConf) (map[string]struct{}, error) {
	result := map[string]struct{}{
		"Added":      {},
		"Changed":    {},
		"Deprecated": {},
		"Fixed":      {},
		"Removed":    {},
		"Security":   {},
	}
	if args == nil {
		return result, nil
	}

	/*
		for _, arg := range args {
			allowed, ok := arg.(string)
			if !ok {
				return nil, fmt.Errorf("expected string, got %T", arg)
			}
			result[allowed] = struct{}{}
		}
	*/
	return result, nil
}
