package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type SubsectionNamming struct{}

func (r SubsectionNamming) Apply(changes model.Changelog, failures chan linting.Failure, args linting.RuleArgs) {
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

func (r SubsectionNamming) allowedSubsections(args linting.RuleArgs) (map[string]struct{}, error) {
	result := map[string]struct{}{
		"Added":      {},
		"Changed":    {},
		"Deprecated": {},
		"Fixed":      {},
		"Removed":    {},
		"Security":   {},
	}

	allowed, _ := args.([]any)
	if len(allowed) == 0 { // no conf provided
		return result, nil
	}

	result = make(map[string]struct{}, len(allowed))
	for _, a := range allowed {
		allow, ok := a.(string)
		if !ok {
			return nil, fmt.Errorf("expected %v to be a string, got (GO)type %T)", a, a)
		}
		result[allow] = struct{}{}
	}

	return result, nil
}
