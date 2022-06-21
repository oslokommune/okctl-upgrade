package upgrade

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func mockCluster() v1alpha1.Cluster {
	return v1alpha1.Cluster{
		Metadata: v1alpha1.ClusterMeta{Name: "mock-cluster"},
		Github:   v1alpha1.ClusterGithub{OutputPath: "infrastructure"},
	}
}

func contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

const appTemplate = `apiVersion: argoproj.io/v1alpha1
kind: Application
`

func TestScanForRelevantApps(t *testing.T) {
	testCases := []struct {
		name        string
		withCluster v1alpha1.Cluster
		withFs      *afero.Afero
		expectApps  []string
	}{
		{
			name:        "Should return correct apps",
			withCluster: mockCluster(),
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				_ = fs.MkdirAll("/infrastructure/mock-cluster/argocd/applications", defaultFolderPermissions)
				_ = fs.WriteReader("/infrastructure/mock-cluster/argocd/applications/mock-app-one.yaml", strings.NewReader(appTemplate))

				return fs
			}(),
			expectApps: []string{"mock-app-one"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			apps, err := scanForRelevantApplications(tc.withFs, tc.withCluster, "/")
			assert.NoError(t, err)

			assert.Equal(t, len(tc.expectApps), len(apps))

			for _, item := range apps {
				assert.True(t, contains(tc.expectApps, item))
			}

			for _, item := range tc.expectApps {
				assert.True(t, contains(apps, item))
			}
		})
	}
}

const namespaceTemplate = `apiVersion: v1
kind: Namespace
metadata:
  name: {{.Name}}`

func namespace(name string) io.Reader {
	t := template.Must(template.New("namespace").Parse(namespaceTemplate))

	buf := bytes.Buffer{}

	_ = t.Execute(&buf, struct {
		Name string
	}{Name: name})

	return &buf
}

func addAppNamespace(t *testing.T, fs *afero.Afero, appName string, namespaceName string) {
	appBaseDir := path.Join("/infrastructure/applications", appName, "base")

	err := fs.MkdirAll(appBaseDir, defaultFolderPermissions)
	assert.NoError(t, err)

	err = fs.WriteReader(path.Join(appBaseDir, "namespace.yaml"), namespace(namespaceName))
	assert.NoError(t, err)
}

func TestMigrateApplication(t *testing.T) {
	testCases := []struct {
		name             string
		withFs           *afero.Afero
		withCluster      v1alpha1.Cluster
		withAppName      string
		expectNamespaces []string
	}{
		{
			name: "Should work",
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				addAppNamespace(t, fs, "mock-app-one", "mock-namespace")

				return fs
			}(),
			withCluster:      mockCluster(),
			withAppName:      "mock-app-one",
			expectNamespaces: []string{"mock-namespace"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := migrateApplication(tc.withFs, tc.withCluster, "/", tc.withAppName)
			assert.NoError(t, err)

			for _, ns := range tc.expectNamespaces {
				exists, err := tc.withFs.Exists(path.Join(
					"/",
					tc.withCluster.Github.OutputPath,
					tc.withCluster.Metadata.Name,
					"argocd",
					"namespaces",
					fmt.Sprintf("%s.yaml", ns),
				))
				assert.NoError(t, err)

				assert.True(t, exists)
			}
		})
	}
}
