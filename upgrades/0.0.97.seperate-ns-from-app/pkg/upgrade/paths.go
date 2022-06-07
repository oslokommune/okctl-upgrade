package upgrade

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
)

const (
	argocdConfigDir           = "argocd"
	argocdConfigNamespacesDir = "namespaces"
	defaultFolderPermissions  = 0o600
	defaultReadmeFilename     = "README.md"
)

func argoCDConfigDir(cluster v1alpha1.Cluster) string {
	return path.Join(cluster.Github.OutputPath, cluster.Metadata.Name, argocdConfigDir)
}

func namespacesDir(cluster v1alpha1.Cluster) string {
	return path.Join(argoCDConfigDir(cluster), argocdConfigNamespacesDir)
}

// getRepositoryRootDirectory returns the absolute path of the repository root no matter what the current working
// directory of the repository the user is in
func getRepositoryRootDirectory() (string, error) {
	result, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("getting repository root directory: %w", err)
	}

	pathAsString := string(bytes.Trim(result, "\n"))

	return pathAsString, nil
}
