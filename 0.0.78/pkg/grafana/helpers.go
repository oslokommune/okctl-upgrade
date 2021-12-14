package grafana

import "fmt"

func validateVersion(expected string, actual string) error {
	if expected != actual {
		return fmt.Errorf("expected %s, got %s", expected, actual)
	}

	return nil
}
