package eksctl

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
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

func GenerateKubeconfig(fs *afero.Afero, clusterName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("acquiring home dir: %w", err)
	}

	eksctlPath, err := acquireEksctlPath(fs, os.UserHomeDir)
	if err != nil {
		return "", fmt.Errorf("acquiring eksctl path: %w", err)
	}

	err = runEksctlCommand(eksctlPath, "utils", "write-kubeconfig", "--auto-kubeconfig", "-c", clusterName)
	if err != nil {
		return "", fmt.Errorf("running command: %w", err)
	}

	return path.Join(homeDir, ".kube", "eksctl", "clusters", clusterName), nil
}
