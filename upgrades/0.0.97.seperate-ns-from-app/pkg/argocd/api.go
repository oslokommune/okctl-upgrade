package argocd

import (
	_ "embed"
	"fmt"
	"path"
	"strings"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/paths"
	"github.com/spf13/afero"
)

func EnableNamespacesSync(logger debugLogger, dryRun bool, fs *afero.Afero, cluster v1alpha1.Cluster) error {
	if !dryRun {
		logger.Debug("Preparing directory structure")

		err := fs.MkdirAll(paths.RelativeNamespacesDir(cluster), paths.DefaultFolderPermissions)
		if err != nil {
			return fmt.Errorf("preparing namespaces dir: %w", err)
		}

		err = fs.WriteReader(
			path.Join(paths.RelativeNamespacesDir(cluster), paths.DefaultReadmeFilename),
			strings.NewReader(namespacesReadmeTemplate),
		)
		if err != nil {
			return fmt.Errorf("creating namespaces readme: %w", err)
		}
	}

	logger.Debug("Adding namespaces ArgoCD application")

	argoApp, err := scaffoldApplication(cluster, "namespaces", paths.RelativeNamespacesDir(cluster))
	if err != nil {
		return fmt.Errorf("scaffolding ArgoCD application: %w", err)
	}

	if !dryRun {
		err = fs.WriteReader(path.Join(paths.RelativeArgoCDConfigDir(cluster), "namespaces.yaml"), argoApp)
		if err != nil {
			return fmt.Errorf("writing ArgoCD application: %w", err)
		}
	}

	return nil
}

//go:embed templates/namespaces-readme.md
var namespacesReadmeTemplate string
