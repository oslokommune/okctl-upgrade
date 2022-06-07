package upgrade

import (
	"context"
	"fmt"
	"path"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/logging"
)

func Start(ctx context.Context, logger logging.Logger, fs *afero.Afero, flags cmdflags.Flags, cluster v1alpha1.Cluster) error {
	absoluteRepositoryRootDir, err := getRepositoryRootDirectory()
	if err != nil {
		return fmt.Errorf("acquiring repository root dir: %w", err)
	}

	absoluteNamespacesDir := path.Join(absoluteRepositoryRootDir, namespacesDir(cluster))

	err = setupNamespacesSync(fs, cluster)
	if err != nil {
		return fmt.Errorf("adding namespaces app manifest: %w", err)
	}

	err = migrateExistingNamespacesToDir(fs, cluster, absoluteNamespacesDir)
	if err != nil {
		return fmt.Errorf("migrating existing namespaces: %w", err)
	}

	return nil
}
