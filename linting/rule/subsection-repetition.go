package rule

import (
	"fmt"

	"github.com/chavacava/changelogger/linting"
	"github.com/chavacava/changelogger/model"
)

type SubsectionRepetition struct{}

func (r SubsectionRepetition) Apply(changes model.Changelog, failures chan linting.Failure, _ []any) {
	for _, version := range changes.Versions {
		seen := map[string]struct{}{}
		for _, subsection := range version.Subsections {
			name := subsection.Name
			_, alreadySeen := seen[name]
			if alreadySeen {
				msg := fmt.Sprintf("duplicated subsection %q", name)
				failures <- linting.Failure{Originator: r.Name(), Message: msg}
			}
			seen[name] = struct{}{}
		}
	}
}

func (SubsectionRepetition) Name() string {
	return "subsection-repetition"
}
