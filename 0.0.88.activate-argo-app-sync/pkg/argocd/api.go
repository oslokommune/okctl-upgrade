package argocd

import (
	"bytes"
	"fmt"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/scaffold"
	"github.com/oslokommune/okctl/pkg/scaffold/resources"
	"github.com/spf13/afero"
	"path"
)

// SetupApplicationsSync knows how to get ArgoCD to automatically synchronize a folder
func SetupApplicationsSync(opts SetupApplicationsSyncOpts) error {
	app := v1alpha1.Application{Metadata: v1alpha1.ApplicationMeta{
		Name:      "cluster-applications",
		Namespace: "argocd",
	}}

	relativeArgoCDManifestPath := path.Join(
		getArgoCDClusterConfigDir(opts.Cluster),
		defaultArgoCDSyncApplicationsManifestName,
	)
	relativeApplicationsSyncDir := path.Join(
		getArgoCDClusterConfigDir(opts.Cluster),
		defaultApplicationsSyncDirName,
	)

	argoCDApplication := resources.CreateArgoApp(app, opts.Cluster.Github.URL(), relativeApplicationsSyncDir)

	rawArgoCDApplication, err := scaffold.ResourceAsBytes(argoCDApplication)
	if err != nil {
		return fmt.Errorf("marshalling ArgoCD application manifest: %w", err)
	}

	err = opts.Fs.WriteReader(relativeArgoCDManifestPath, bytes.NewReader(rawArgoCDApplication))
	if err != nil {
		return fmt.Errorf("writing ArgoCD application manifest: %w", err)
	}

	err = opts.Kubectl.Apply(bytes.NewReader(rawArgoCDApplication))
	if err != nil {
		return fmt.Errorf("applying ArgoCD application manifest: %w", err)
	}

	return nil
}

// MigrateExistingApplicationManifests knows how to move all existing argocd-application manifests to the new sync
// directory
func MigrateExistingApplicationManifests(filesystem *afero.Afero, cluster v1alpha1.Cluster) error {
	rootAppDir := path.Join(cluster.Github.OutputPath, defaultApplicationsDirName)

	relativeArgoCDClusterConfigDir := getArgoCDClusterConfigDir(cluster)
	relativeAppSyncDir := path.Join(relativeArgoCDClusterConfigDir, defaultApplicationsSyncDirName)

	argoCDApplicationManifestPaths, err := getAllArgoCDApplicationManifests(filesystem, rootAppDir)
	if err != nil {
		return fmt.Errorf("acquiring all ArgoCD application manifest paths: %w", err)
	}

	for _, sourcePath := range argoCDApplicationManifestPaths {
		appName := getApplicationNameFromPath(rootAppDir, sourcePath)

		destinationPath := path.Join(relativeAppSyncDir, fmt.Sprintf("%s.yaml", appName))

		err = filesystem.Rename(sourcePath, destinationPath)
		if err != nil {
			return fmt.Errorf("moving file: %w", err)
		}
	}

	return nil
}
