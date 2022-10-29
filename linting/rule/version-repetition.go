package rule

import (
	"fmt"

	"github.com/chavacava/changelogger/linting"
	"github.com/chavacava/changelogger/model"
)

type VersionRepetition struct{}

func (r VersionRepetition) Apply(changes model.Changelog, failures chan linting.Failure, _ []any) {
	seen := map[string]struct{}{}
	for _, version := range changes.Versions {
		_, alreadySeen := seen[version.Version]
		if alreadySeen {
			msg := fmt.Sprintf("duplicated version %q", version.Version)
			failures <- linting.Failure{Originator: r.Name(), Message: msg}
		}
		seen[version.Version] = struct{}{}
	}
}

func (VersionRepetition) Name() string {
	return "version-repetition"
}
