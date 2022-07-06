package manifest

import (
	"fmt"
	"github.com/oslokommune/okctl-upgrade/upgrades/bump-eks/pkg/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
)

func Cluster(fs *afero.Afero, clusterManifestPath string) (v1alpha1.Cluster, error) {
	rawManifest, err := fs.ReadFile(clusterManifestPath)
	if err != nil {
		return v1alpha1.Cluster{}, fmt.Errorf("reading manifest: %w", err)
	}

	clusterManifest, err := v1alpha1.Parse(rawManifest)
	if err != nil {
		return v1alpha1.Cluster{}, fmt.Errorf("parsing manifest: %w", err)
	}

	return clusterManifest, nil
}
