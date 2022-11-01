package rule

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chavacava/changelog-lint/linting"
	"github.com/chavacava/changelog-lint/model"
)

type VersionOrder struct{}

func (r VersionOrder) Apply(changes model.Changelog, failures chan linting.Failure, _ []any) {
	previousVersion := ""
	for _, version := range changes.Versions {
		if previousVersion != "" && version.Version == "Unreleased" {
			msg := "version Unreleased must be at the top of the version list"
			failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: version.Position}
			continue
		}

		if previousVersion != "" && r.compareVersions(previousVersion, version.Version) < 0 {
			msg := fmt.Sprintf("version %s is not well sorted", version.Version)
			failures <- linting.Failure{RuleName: r.Name(), Message: msg, Position: version.Position}
		}
		previousVersion = version.Version
	}
}

func (VersionOrder) Name() string {
	return "version-order"
}

func (VersionOrder) compareVersions(v1 string, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")
	if len(v1Parts) != 3 || len(v2Parts) != 3 {
		panic(fmt.Sprintf("%q or %q is not a semver string", v1, v2))
	}

	for i, p1 := range v1Parts {
		p2 := v2Parts[i]
		value1, err := strconv.Atoi(p1)
		if err != nil {
			panic(fmt.Sprintf("%q is not a semver string: %v", v1, err))
		}
		value2, err := strconv.Atoi(p2)
		if err != nil {
			panic(fmt.Sprintf("%q is not a semver string: %v", v2, err))
		}

		if value1 > value2 {
			return 1
		}

		if value1 < value2 {
			return -1
		}

		continue
	}
	return 0
}
