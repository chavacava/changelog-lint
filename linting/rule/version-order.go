package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type VersionOrder struct{}

func (r VersionOrder) Apply(changes model.Changelog, failures chan linting.Failure, _ []any) {
	previousVersion := ""
	for _, version := range changes.Versions {
		if previousVersion != "" && version.Version > previousVersion {
			msg := fmt.Sprintf("version %q is not well sorted", version.Version)
			failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: version.Position}
		}
		previousVersion = version.Version
	}
}

func (VersionOrder) Name() string {
	return "version-order"
}
