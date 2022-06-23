package upgrade

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
)

func getApplicationsInCluster(fs *afero.Afero, cluster v1alpha1.Cluster, absoluteRepositoryRoot string) ([]string, error) {
	absoluteApplicationsDir := path.Join(absoluteRepositoryRoot, argoCDApplicationsDir(cluster))

	files, err := fs.ReadDir(absoluteApplicationsDir)
	if err != nil {
		return nil, fmt.Errorf("retrieving items in app dir: %w", err)
	}

	apps := make([]string, 0)

	for _, potentialApp := range files {
		relevant, err := isPathAnArgoCDApplication(fs, path.Join(absoluteApplicationsDir, potentialApp.Name()))
		if err != nil {
			return nil, fmt.Errorf("checking relevance: %w", err)
		}

		if relevant {
			apps = append(apps, filenameWithoutExtension(potentialApp.Name()))
		}
	}

	return apps, nil
}
