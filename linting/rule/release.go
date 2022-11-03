package rule

import (
	"fmt"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type Release struct{}

func (r Release) Apply(changes model.Changelog, failures chan linting.Failure, args linting.RuleArgs) {
	if len(changes.Versions) < 1 {
		msg := "at least one version required in release mode"
		failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: 0}
		return
	}

	wantVersion, ok := args.(string)
	if !ok {
		msg := fmt.Sprintf("expected release version argument to be a string, got %T instead", args)
		failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: 0}
		return
	}

	headVersion := changes.Versions[0]
	gotVersion := headVersion.Version
	if gotVersion != wantVersion {
		msg := fmt.Sprintf("expected release version to be %s, got %s instead", wantVersion, gotVersion)
		failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: headVersion.Position}
	}

	for _, version := range changes.Versions {
		if version.Version == "Unreleased" {
			msg := "version Unreleased forbidden in release mode"
			failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: version.Position}
			break
		}
	}
}

func (Release) Name() string {
	return "release"
}
