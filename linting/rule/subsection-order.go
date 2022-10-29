package rule

import (
	"fmt"

	"github.com/chavacava/changelogger/linting"
	"github.com/chavacava/changelogger/model"
)

type SubsectionOrder struct{}

func (r SubsectionOrder) Apply(changes model.Changelog, failures chan linting.Failure, _ []any) {
	for _, version := range changes.Versions {
		previousName := ""
		for _, subsection := range version.Subsections {
			if previousName != "" && subsection.Name < previousName {
				msg := fmt.Sprintf("subsection %q is not sorted alphabetically in version %v", subsection.Name, version.Version)
				failures <- linting.Failure{Originator: r.Name(), Message: msg}
			}
			previousName = subsection.Name
		}
	}
}

func (SubsectionOrder) Name() string {
	return "subsection-order"
}
