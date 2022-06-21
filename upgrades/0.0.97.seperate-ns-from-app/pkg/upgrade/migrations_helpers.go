package upgrade

import (
	"fmt"
	"path"
	"strings"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"

	"sigs.k8s.io/yaml"
)

func copyFile(fs *afero.Afero, sourcePath string, destinationPath string) error {
	content, err := fs.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("reading: %w", err)
	}

	defer func() {
		_ = content.Close()
	}()

	err = fs.WriteReader(destinationPath, content)
	if err != nil {
		return fmt.Errorf("writing: %w", err)
	}

	return nil
}

func acquireNamespaceName(fs *afero.Afero, namespacePath string) (string, error) {
	content, err := fs.ReadFile(namespacePath)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	var ns namespaceManifest

	err = yaml.Unmarshal(content, &ns)
	if err != nil {
		return "", fmt.Errorf("unmarshalling: %w", err)
	}

	return ns.Metadata.Name, nil
}

func isRelevant(fs *afero.Afero, absoluteApplicationPath string) (bool, error) {
	stat, err := fs.Stat(absoluteApplicationPath)
	if err != nil {
		return false, fmt.Errorf("stating file: %w", err)
	}

	if stat.IsDir() {
		return false, nil
	}

	rawFile, err := fs.ReadFile(absoluteApplicationPath)
	if err != nil {
		return false, fmt.Errorf("reading file: %w", err)
	}

	var potentialApp argoCDApplication

	err = yaml.Unmarshal(rawFile, &potentialApp)
	if err != nil {
		return false, fmt.Errorf("unmarshalling: %w", err)
	}

	return potentialApp.Valid(), nil
}

func filenameWithoutExtension(filename string) string {
	return strings.Replace(path.Base(filename), path.Ext(filename), "", 1)
}

func scanForRelevantApplications(fs *afero.Afero, cluster v1alpha1.Cluster, absoluteRepositoryRoot string) ([]string, error) {
	absoluteApplicationsDir := path.Join(absoluteRepositoryRoot, argoCDApplicationsDir(cluster))

	files, err := fs.ReadDir(absoluteApplicationsDir)
	if err != nil {
		return nil, fmt.Errorf("retrieving items in app dir: %w", err)
	}

	apps := make([]string, 0)

	for _, potentialApp := range files {
		relevant, err := isRelevant(fs, path.Join(absoluteApplicationsDir, potentialApp.Name()))
		if err != nil {
			return nil, fmt.Errorf("checking relevance: %w", err)
		}

		if relevant {
			apps = append(apps, filenameWithoutExtension(potentialApp.Name()))
		}
	}

	return apps, nil
}
