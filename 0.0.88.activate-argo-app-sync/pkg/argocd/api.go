package argocd

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
)

// SetupApplicationsSync knows how to get ArgoCD to automatically synchronize a folder
func SetupApplicationsSync(opts SetupApplicationsSyncOpts) error {
	log := opts.Logger

	log.Info("Setting up application synchronization")

	relativeArgoCDManifestPath := path.Join(
		getArgoCDClusterConfigDir(opts.Cluster),
		defaultArgoCDSyncApplicationsManifestName,
	)
	relativeApplicationsSyncDir := path.Join(
		getArgoCDClusterConfigDir(opts.Cluster),
		defaultApplicationsSyncDirName,
	)

	log.Infof("Creating new application sync directory %s", relativeApplicationsSyncDir)

	err := createDirectory(opts.Fs, opts.DryRun, relativeApplicationsSyncDir)
	if err != nil {
		return fmt.Errorf("creating applications sync directory: %w", err)
	}

	log.Info("Installing ArgoCD application for application sync directory")

	err = installArgoCDApplicationForSyncDirectory(installArgoCDApplicationForSyncDirectoryOpts{
		DryRun:                        opts.DryRun,
		Fs:                            opts.Fs,
		Kubectl:                       opts.Kubectl,
		IACRepoURL:                    opts.Cluster.Github.URL(),
		ApplicationsSyncDir:           relativeApplicationsSyncDir,
		ArgoCDApplicationManifestPath: relativeArgoCDManifestPath,
	})
	if err != nil {
		return fmt.Errorf("installing ArgoCD application: %w", err)
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
