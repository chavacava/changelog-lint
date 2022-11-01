package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type SubsectionRepetition struct{}

func (r SubsectionRepetition) Apply(changes model.Changelog, failures chan linting.Failure, _ linting.RuleArgs) {
	for _, version := range changes.Versions {
		seen := map[string]struct{}{}
		for _, subsection := range version.Subsections {
			name := subsection.Name
			_, alreadySeen := seen[name]
			if alreadySeen {
				msg := fmt.Sprintf("duplicated subsection %q", name)
				failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: subsection.Position}
			}
			seen[name] = struct{}{}
		}
	}
}

func (SubsectionRepetition) Name() string {
	return "subsection-repetition"
}
