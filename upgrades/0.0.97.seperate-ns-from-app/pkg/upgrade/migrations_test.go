package upgrade

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/paths"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.97.seperate-ns-from-app/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func mockCluster(name string) v1alpha1.Cluster {
	return v1alpha1.Cluster{
		Metadata: v1alpha1.ClusterMeta{Name: name},
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
			withCluster: mockCluster("mock-cluster"),
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				_ = fs.MkdirAll("/infrastructure/mock-cluster/argocd/applications", paths.DefaultFolderPermissions)
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

			apps, err := getApplicationsInCluster(tc.withFs, tc.withCluster, "/")
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

func addOldAppNamespace(t *testing.T, fs *afero.Afero, appName string, namespaceName string) {
	appBaseDir := path.Join("/infrastructure/applications", appName, "base")

	err := fs.MkdirAll(appBaseDir, paths.DefaultFolderPermissions)
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

				addOldAppNamespace(t, fs, "mock-app-one", "mock-namespace")

				return fs
			}(),
			withCluster:      mockCluster("mock-cluster"),
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

// One app "exists" when the /infrastructure/applications/<app name> directory exists
func createApp(t *testing.T, fs *afero.Afero, appName string) {
	absOutputDir := path.Join("/", "infrastructure")
	absAppDir := path.Join(absOutputDir, "applications", appName)
	absBaseDir := path.Join(absAppDir, "base")

	err := fs.MkdirAll(absBaseDir, paths.DefaultFolderPermissions)
	assert.NoError(t, err)
}

// One cluster "exists" when the /infrastructure/<cluster name> directory exists
func createCluster(t *testing.T, fs *afero.Afero, clusterName string) {
	absOutputDir := path.Join("/", "infrastructure")
	absArgoCDApplicationsConfigDir := path.Join(absOutputDir, clusterName, "argocd", "applications")

	err := fs.MkdirAll(absArgoCDApplicationsConfigDir, paths.DefaultFolderPermissions)
	assert.NoError(t, err)
}

const argoCDApplicationTemplate = `apiVersion: argoproj.io/v1alpha1
kind: Application
`

// An app is added to a cluster when the /infrastructure/<cluster name>/argocd/applications/<app name>.yaml file exists
func addAppToCluster(t *testing.T, fs *afero.Afero, appName string, clusterName string) {
	absOutputDir := path.Join("/", "infrastructure")
	absArgoCDApplicationsConfigDir := path.Join(absOutputDir, clusterName, "argocd", "applications")

	err := fs.MkdirAll(absArgoCDApplicationsConfigDir, paths.DefaultFolderPermissions)
	assert.NoError(t, err)

	err = fs.WriteReader(
		path.Join(absArgoCDApplicationsConfigDir, fmt.Sprintf("%s.yaml", appName)),
		strings.NewReader(argoCDApplicationTemplate),
	)
	assert.NoError(t, err)
}

func addNewAppNamespace(t *testing.T, fs *afero.Afero, clusterName string, namespaceName string) {
	absOutputDir := path.Join("/", "infrastructure")
	absArgoCDNamespacesConfigDir := path.Join(absOutputDir, clusterName, "argocd", "namespaces")

	err := fs.MkdirAll(absArgoCDNamespacesConfigDir, paths.DefaultFolderPermissions)
	assert.NoError(t, err)

	err = fs.WriteReader(
		path.Join(absArgoCDNamespacesConfigDir, fmt.Sprintf("%s.yaml", namespaceName)),
		namespace(namespaceName),
	)
	assert.NoError(t, err)
}

func TestRemoveRedundantNamespacesFromBase(t *testing.T) {
	testCases := []struct {
		name                     string
		withFs                   *afero.Afero
		expectedNonExistantPaths []string
		expectedExistantPaths    []string
		withCurrentCluster       v1alpha1.Cluster
	}{
		{
			name:               "Should remove base ns with one app and one upgraded cluster",
			withCurrentCluster: mockCluster("mock-prod"),
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				clusterOne := "mock-prod"
				createCluster(t, fs, clusterOne)

				appOne := "mock-app-one"
				createApp(t, fs, appOne)

				appOneNamespace := "apps"
				addOldAppNamespace(t, fs, appOne, appOneNamespace)
				addNewAppNamespace(t, fs, clusterOne, appOneNamespace)

				addAppToCluster(t, fs, appOne, clusterOne)

				return fs
			}(),
			expectedNonExistantPaths: []string{"/infrastructure/applications/mock-app-one/base/namespace.yaml"},
			expectedExistantPaths:    []string{"/infrastructure/mock-prod/argocd/namespaces/apps.yaml"},
		},
		{
			name:               "Should leave ns in base with one upgraded cluster and one not upgraded cluster",
			withCurrentCluster: mockCluster("mock-prod"),
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				clusterOne := "mock-prod"
				createCluster(t, fs, clusterOne)

				clusterTwo := "mock-test"
				createCluster(t, fs, clusterTwo)

				appOne := "mock-app-one"
				createApp(t, fs, appOne)

				appOneNamespace := "apps"
				addOldAppNamespace(t, fs, appOne, appOneNamespace)
				addNewAppNamespace(t, fs, clusterOne, appOneNamespace)

				addAppToCluster(t, fs, appOne, clusterOne)
				addAppToCluster(t, fs, appOne, clusterTwo)

				return fs
			}(),
			expectedNonExistantPaths: []string{},
			expectedExistantPaths: []string{
				"/infrastructure/mock-prod/argocd/namespaces/apps.yaml",
				"/infrastructure/applications/mock-app-one/base/namespace.yaml",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := removeRedundantNamespacesFromBase(tc.withFs, tc.withCurrentCluster, "/")
			assert.NoError(t, err)

			for _, currentPath := range tc.expectedExistantPaths {
				exists, err := tc.withFs.Exists(currentPath)
				assert.NoError(t, err)

				assert.True(t, exists)
			}

			for _, currentPath := range tc.expectedNonExistantPaths {
				exists, err := tc.withFs.Exists(currentPath)
				assert.NoError(t, err)

				assert.False(t, exists)
			}
		})
	}
}

func TestGetApplicationsInCluster(t *testing.T) {
	testCases := []struct {
		name              string
		withFs            *afero.Afero
		withClusterName   string
		expectedAppsFound []string
	}{
		{
			name:            "Should work",
			withClusterName: "mock-cluster",
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				createCluster(t, fs, "mock-cluster")
				createApp(t, fs, "mock-app-one")

				addAppToCluster(t, fs, "mock-app-one", "mock-cluster")

				return fs
			}(),
			expectedAppsFound: []string{"mock-app-one"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			apps, err := getApplicationsInCluster(
				tc.withFs,
				mockCluster(tc.withClusterName),
				"/",
			)
			assert.NoError(t, err)

			assert.Equal(t, len(tc.expectedAppsFound), len(apps))

			for _, app := range apps {
				assert.True(t, contains(tc.expectedAppsFound, app))
			}

			for _, app := range tc.expectedAppsFound {
				assert.True(t, contains(apps, app))
			}
		})
	}
}
