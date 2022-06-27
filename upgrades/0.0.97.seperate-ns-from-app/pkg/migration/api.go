package migration

import (
	"fmt"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
)

func MigrateExistingApplicationNamespacesToCluster(fs *afero.Afero, cluster v1alpha1.Cluster, absoluteRepositoryRoot string) error {
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

	err = removeRedundantNamespacesFromBase(fs, cluster, absoluteRepositoryRoot)
	if err != nil {
		return fmt.Errorf("removing redundant namespace manifests from application base folders: %w", err)
	}

	return nil
}
