package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type SubsectionEmpty struct{}

func (r SubsectionEmpty) Apply(changes model.Changelog, failures chan linting.Failure, _ []any) {
	for _, version := range changes.Versions {
		for _, subsection := range version.Subsections {
			if len(subsection.History) == 0 {
				msg := fmt.Sprintf("empty subsection %q in version %v", subsection.Name, version.Version)
				failures <- linting.Failure{Originator: r.Name(), Message: msg}
			}
		}
	}
}

func (SubsectionEmpty) Name() string {
	return "subsection-empty"
}
