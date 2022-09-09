package main

import (
	"testing"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestGenerateNodegroupNames(t *testing.T) {
	testCases := []struct {
		name                 string
		withCluster          v1alpha1.Cluster
		withClusterVersion   string
		expectNodegroupNames []string
	}{
		{
			name: "Should work",
			withCluster: v1alpha1.Cluster{
				Metadata: v1alpha1.ClusterMeta{
					Name:   "mock-cluster",
					Region: "eu-west-1",
				},
			},
			withClusterVersion: "1.21",
			expectNodegroupNames: []string{
				"ng-generic-1-21-1a",
				"ng-generic-1-21-1b",
				"ng-generic-1-21-1c",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			names := generateNodegroupNames(tc.withCluster.Metadata.Region, tc.withClusterVersion)

			assert.Equal(t, tc.expectNodegroupNames, names)
		})
	}
}
