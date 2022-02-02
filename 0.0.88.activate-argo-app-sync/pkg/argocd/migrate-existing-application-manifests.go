package argocd

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/spf13/afero"
)

// getAllArgoCDApplicationManifests walks rootDir recursively and returns all files named 'argocd-application.yaml'
func getAllArgoCDApplicationManifests(filesystem *afero.Afero, rootDir string) ([]string, error) {
	result := make([]string, 0)

	err := filesystem.Walk(rootDir, func(currentPath string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if info.Name() != defaultArgoCDApplicationManifestName {
			return nil
		}

		result = append(result, path.Join(rootDir, currentPath))

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("gathering ArgoCD application manifests: %w", err)
	}

	return result, nil
}

// getApplicationNameFromPath knows how to extract the application name from the path of an ArgoCD application manifest
func getApplicationNameFromPath(applicationsRootDirectory string, targetPath string) string {
	cleanedPath := strings.Replace(targetPath, applicationsRootDirectory, "", 1)
	cleanedPath = strings.TrimPrefix(cleanedPath, "/")

	parts := strings.Split(cleanedPath, "/")

	return parts[0]
}
