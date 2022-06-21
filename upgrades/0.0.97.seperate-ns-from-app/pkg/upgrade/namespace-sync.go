package upgrade

import (
	"fmt"
	"path"
	"strings"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/argocd"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
)

func setupNamespacesSync(fs *afero.Afero, cluster v1alpha1.Cluster) error {
	err := fs.MkdirAll(namespacesDir(cluster), defaultFolderPermissions)
	if err != nil {
		return fmt.Errorf("preparing namespaces dir: %w", err)
	}

	err = fs.WriteReader(
		path.Join(namespacesDir(cluster), defaultReadmeFilename),
		strings.NewReader(namespacesReadmeTemplate),
	)
	if err != nil {
		return fmt.Errorf("creating namespaces readme: %w", err)
	}

	argoApp, err := argocd.ScaffoldApplication(cluster, "namespaces", namespacesDir(cluster))
	if err != nil {
		return fmt.Errorf("scaffolding ArgoCD application: %w", err)
	}

	err = fs.WriteReader(path.Join(argoCDConfigDir(cluster), "namespaces.yaml"), argoApp)
	if err != nil {
		return fmt.Errorf("writing ArgoCD application: %w", err)
	}

	return nil
}

const namespacesReadmeTemplate = `# Namespaces/

To remove ownership of a namespace from an individual application, we've set up this folder and made ArgoCD
automatically track changes to it. This folder contains all namespace manifests. Adding manifests to this folder will
automatically apply them to your cluster.

When running the apply application command, okctl will place a namespace manifest here.
`
