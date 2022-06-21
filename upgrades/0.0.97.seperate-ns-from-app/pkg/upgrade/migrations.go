package upgrade

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
)

func migrateExistingNamespacesToDir(fs *afero.Afero, cluster v1alpha1.Cluster, absoluteRepositoryRoot string) error {
	apps, err := scanForRelevantApplications(fs, cluster, absoluteRepositoryRoot)
	if err != nil {
		return fmt.Errorf("scanning for : %w", err)
	}

	for _, app := range apps {
		err = migrateApplication(fs, cluster, absoluteRepositoryRoot, app)
		if err != nil {
			return fmt.Errorf("migrating %s: %w", app, err)
		}
	}

	return nil
}

func migrateApplication(fs *afero.Afero, cluster v1alpha1.Cluster, absoluteRepositoryRoot string, appName string) error {
	absoluteNamespacesDir := path.Join(absoluteRepositoryRoot, namespacesDir(cluster))
	absoluteApplicationBaseDir := path.Join(
		absoluteRepositoryRoot,
		cluster.Github.OutputPath,
		applicationsDir,
		appName,
		applicationBaseDir,
	)

	sourcePath := path.Join(absoluteApplicationBaseDir, "namespace.yaml")

	namespaceName, err := acquireNamespaceName(fs, sourcePath)
	if err != nil {
		return fmt.Errorf("acquiring namespace name: %w", err)
	}

	destinationPath := path.Join(absoluteNamespacesDir, fmt.Sprintf("%s.yaml", namespaceName))

	err = copyFile(fs, sourcePath, destinationPath)
	if err != nil {
		return fmt.Errorf("copying: %w", err)
	}

	return nil
}
