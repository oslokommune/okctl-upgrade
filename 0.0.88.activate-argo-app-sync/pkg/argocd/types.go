package argocd

import (
	"github.com/oslokommune/okctl-upgrade/0.0.88.activate-argo-app-sync/pkg/kubectl"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
)

const (
	// defaultApplicationsDirName defines the name of the directory that contains all the okctl applications
	defaultApplicationsDirName = "applications"
	// defaultApplicationsSyncDirName defines the name of the directory that gets automatically synchronized
	defaultApplicationsSyncDirName = "applications"
	// defaultArgoCDSyncApplicationsManifestName defines the name of the ArgoCD application manifest that enables
	// synchronization of the defaultApplicationsSyncDirName
	defaultArgoCDSyncApplicationsManifestName = "cluster-applications.yaml"
	// defaultArgoCDApplicationManifestName defines the "old" name of the individual application ArgoCD app manifests
	defaultArgoCDApplicationManifestName = "argocd-application.yaml"
	// defaultArgoCDClusterConfigDirName defines the name of the cluster specific ArgoCD configuration directory
	defaultArgoCDClusterConfigDirName = "argocd"
	// defaultFolderPermissions defines the default permissions for the ArgoCD config directory and applications sync
	// directory
	defaultFolderPermissions = 0o700
)

// SetupApplicationsSyncOpts defines necessary data required to setup application synchronization
type SetupApplicationsSyncOpts struct {
	Fs      *afero.Afero
	Cluster v1alpha1.Cluster
	Kubectl kubectl.Client
}
