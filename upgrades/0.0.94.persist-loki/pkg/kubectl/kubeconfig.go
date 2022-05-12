package kubectl

import (
	"fmt"
	"os"
	"path"
)

func acquireKubeconfigPath(clusterName string) (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("acquiring user home directory: %w", err)
	}

	return path.Join(
		userDir,
		defaultOkctlConfigDirName,
		defaultCredentialsDirName,
		clusterName,
		defaultKubeconfigFilename,
	), nil
}
