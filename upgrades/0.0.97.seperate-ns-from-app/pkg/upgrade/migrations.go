package upgrade

import (
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
)

func migrateExistingNamespacesToDir(fs *afero.Afero, cluster v1alpha1.Cluster, absoluteNamespacesDir string) error {
	return nil
}
