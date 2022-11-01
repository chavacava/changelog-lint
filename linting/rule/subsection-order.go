package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type SubsectionOrder struct{}

func (r SubsectionOrder) Apply(changes model.Changelog, failures chan linting.Failure, _ linting.RuleConf) {
	for _, version := range changes.Versions {
		previousName := ""
		for _, subsection := range version.Subsections {
			if previousName != "" && subsection.Name < previousName {
				msg := fmt.Sprintf("subsection %q is not sorted alphabetically in version %v", subsection.Name, version.Version)
				failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: subsection.Position}
			}
			previousName = subsection.Name
		}
	}
}

func (SubsectionOrder) Name() string {
	return "subsection-order"
}
