package eksctl

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"strings"
)

func CreateServiceUser(fs *afero.Afero, clusterName string, name string, policies []string) error {
	eksctlPath, err := acquireEksctlPath(fs, os.UserHomeDir)
	if err != nil {
		return fmt.Errorf("acquiring eksctl path: %w", err)
	}

	err = runEksctlCommand(eksctlPath,
		"create", "iamserviceaccount",
		"--name", name,
		"--namespace", defaultMonitoringNamespace,
		"--cluster", clusterName,
		"--role-name", fmt.Sprintf("okctl-%s-loki", clusterName),
		"--attach-policy-arn", strings.Join(policies, ","),
		"--approve",
		"--override-existing-serviceaccounts",
	)
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}
