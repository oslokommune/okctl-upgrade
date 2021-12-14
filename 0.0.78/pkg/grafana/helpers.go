package grafana

import (
	"fmt"

	"github.com/Masterminds/semver"
)

func validateVersion(expected *semver.Version, actual *semver.Version) error {
	if !actual.Equal(expected) {
		return fmt.Errorf("expected %s, got %s", expected, actual)
	}

	return nil
}
