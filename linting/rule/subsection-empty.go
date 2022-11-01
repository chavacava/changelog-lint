package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type SubsectionEmpty struct{}

func (r SubsectionEmpty) Apply(changes model.Changelog, failures chan linting.Failure, _ linting.RuleArgs) {
	for _, version := range changes.Versions {
		for _, subsection := range version.Subsections {
			if len(subsection.History) == 0 {
				msg := fmt.Sprintf("empty subsection %q in version %v", subsection.Name, version.Version)
				failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: subsection.Position}
			}
		}
	}
}

func (SubsectionEmpty) Name() string {
	return "subsection-empty"
}
