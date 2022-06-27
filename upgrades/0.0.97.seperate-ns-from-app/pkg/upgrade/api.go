package upgrade

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/argocd"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/migration"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/paths"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/logging"
)

func Start(_ context.Context, logger logging.Logger, fs *afero.Afero, flags cmdflags.Flags, cluster v1alpha1.Cluster) error {
	absoluteRepositoryRootDir, err := paths.GetRepositoryRootDirectory()
	if err != nil {
		return fmt.Errorf("acquiring repository root dir: %w", err)
	}

	logger.Debug("Enabling namespace synchronization")

	err = argocd.EnableNamespacesSync(&logger, fs, cluster)
	if err != nil {
		return fmt.Errorf("adding namespaces app manifest: %w", err)
	}

	logger.Debug("Migrate application owned namespaces")

	err = migration.MigrateExistingApplicationNamespacesToCluster(logger, fs, cluster, absoluteRepositoryRootDir)
	if err != nil {
		return fmt.Errorf("migrating existing namespaces: %w", err)
	}

	return nil
}
