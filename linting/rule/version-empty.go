package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type VersionEmpty struct{}

func (r VersionEmpty) Apply(changes model.Changelog, failures chan linting.Failure, _ linting.RuleArgs) {
	seen := map[string]struct{}{}
	for _, version := range changes.Versions {
		if len(version.Subsections) == 0 {
			msg := fmt.Sprintf("empty version %q", version.Version)
			failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: version.Position}
		}
		seen[version.Version] = struct{}{}
	}
}

func (VersionEmpty) Name() string {
	return "version-empty"
}
