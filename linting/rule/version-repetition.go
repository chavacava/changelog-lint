package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type VersionRepetition struct{}

func (r VersionRepetition) Apply(changes model.Changelog, failures chan linting.Failure, _ linting.RuleArgs) {
	seen := map[string]struct{}{}
	for _, version := range changes.Versions {
		_, alreadySeen := seen[version.Version]
		if alreadySeen {
			msg := fmt.Sprintf("duplicated version %s", version.Version)
			failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: version.Position}
		}
		seen[version.Version] = struct{}{}
	}
}

func (VersionRepetition) Name() string {
	return "version-repetition"
}
