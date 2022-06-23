package upgrade

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
)

func migrateExistingNamespacesToDir(fs *afero.Afero, cluster v1alpha1.Cluster, absoluteRepositoryRoot string) error {
	apps, err := getApplicationsInCluster(fs, cluster, absoluteRepositoryRoot)
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

	namespaceName, err := getNamespaceName(fs, sourcePath)
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

func isFullyMigrated(fs *afero.Afero, cluster v1alpha1.Cluster, absoluteRepositoryRoot string, appName string) (bool, error) {
	absAppBasePath := path.Join(
		absoluteRepositoryRoot,
		cluster.Github.OutputPath,
		applicationsDir,
		appName,
		applicationBaseDir,
	)

	exists, err := fs.Exists(path.Join(absAppBasePath, "namespace.yaml"))
	if err != nil {
		return false, fmt.Errorf(": %w", err)
	}

	return !exists, nil
}

func removeRedundantNamespacesFromBase(fs *afero.Afero, cluster v1alpha1.Cluster, absoluteRepositoryRoot string) error {
	apps, err := getApplicationsInCluster(fs, cluster, absoluteRepositoryRoot)
	if err != nil {
		return fmt.Errorf("acquiring apps: %w", err)
	}

	for _, app := range apps {
		migrated, err := isFullyMigrated(fs, cluster, absoluteRepositoryRoot, app)
		if err != nil {
			return fmt.Errorf("checking for base namespace: %w", err)
		}

		if migrated {
			continue
		}

		absoluteApplicationDir := path.Join(absoluteRepositoryRoot, cluster.Github.OutputPath, applicationsDir, app)
		absoluteNamespacePath := path.Join(absoluteApplicationDir, applicationBaseDir, "namespace.yaml")

		namespaceName, err := getNamespaceName(fs, absoluteNamespacePath)
		if err != nil {
			return fmt.Errorf("acquiring namespace name: %w", err)
		}

		cleanable, err := allAdjacentClustersHasNamespace(fs, absoluteRepositoryRoot, cluster, namespaceName)
		if err != nil {
			return fmt.Errorf("checking if all adjacent clusters has namespace: %w", err)
		}

		if cleanable {
			err = fs.Remove(absoluteNamespacePath)
			if err != nil {
				return fmt.Errorf("removing base namespace: %w", err)
			}
		}
	}

	return nil
}

func clusterHasNamespace(fs *afero.Afero, absoluteRepositoryOutputDir string, clusterName string, namespaceName string) (bool, error) {
	potentialNamespacePath := path.Join(
		absoluteRepositoryOutputDir,
		clusterName,
		argocdConfigDir,
		"namespaces",
		fmt.Sprintf("%s.yaml", namespaceName),
	)

	exists, err := fs.Exists(potentialNamespacePath)
	if err != nil {
		return false, fmt.Errorf("checking existence: %w", err)
	}

	return exists, nil
}

// Adjacent meaning other clusters in an output folder, i.e.: "infrastructure"
func allAdjacentClustersHasNamespace(fs *afero.Afero, absoluteRepositoryRoot string, cluster v1alpha1.Cluster, namespaceName string) (bool, error) {
	absoluteRepositoryOutputDir := path.Join(absoluteRepositoryRoot, cluster.Github.OutputPath)

	clusters, err := getClusters(fs, absoluteRepositoryOutputDir)
	if err != nil {
		return false, fmt.Errorf("acquiring adjacent clusters: %w", err)
	}

	for _, currentCluster := range clusters {
		hasNamespace, err := clusterHasNamespace(fs, absoluteRepositoryOutputDir, currentCluster, namespaceName)
		if err != nil {
			return false, fmt.Errorf("checking for namespace in cluster: %w", err)
		}

		if !hasNamespace {
			return false, nil
		}
	}

	return true, nil
}

func getClusters(fs *afero.Afero, absoluteRepositoryOutputDir string) ([]string, error) {
	targets, err := fs.ReadDir(absoluteRepositoryOutputDir)
	if err != nil {
		return nil, fmt.Errorf("listing dir: %w", err)
	}

	clusters := make([]string, 0)

	for _, target := range targets {
		name := target.Name()

		if name == applicationsDir {
			continue
		}

		if target.IsDir() {
			clusters = append(clusters, target.Name())
		}
	}

	return clusters, nil
}
