package eksctl

import (
	"io"
	"testing"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/manifest/apis/okctl.io/v1alpha1"
	goldie "github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestGenerateClusterConfig(t *testing.T) {
	testCases := []struct {
		name            string
		withClusterName string
		withAccountID   string
		withRegion      string
		withNodeGroups  []string
	}{
		{
			name:            "Should generate expected clusterconfig with sane values",
			withClusterName: "mock-cluster",
			withAccountID:   "0123456789012",
			withRegion:      "eu-north-1",
			withNodeGroups: []string{
				"eksctl-mock-cluster-nodegroup-ng-generic-1-21-1a-AAAAAAAAAA",
				"eksctl-mock-cluster-nodegroup-ng-generic-1-21-1b-AAAAAAAAAA",
				"eksctl-mock-cluster-nodegroup-ng-generic-1-21-1c-AAAAAAAAAA",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cluster := v1alpha1.NewCluster()
			cluster.Metadata.Name = tc.withClusterName
			cluster.Metadata.AccountID = tc.withAccountID
			cluster.Metadata.Region = tc.withRegion

			cfg, err := GenerateClusterConfig(cluster, tc.withNodeGroups)
			assert.NoError(t, err)

			raw, err := io.ReadAll(cfg)
			assert.NoError(t, err)

			g := goldie.New(t)

			g.Assert(t, tc.name, raw)
		})
	}
}

func TestParseNodeGroups(t *testing.T) {
	testCases := []struct {
		name               string
		withNodegroupNames []string
		expectNodeGroups   []clusterConfigTemplateOptsNodeGroup
	}{
		{
			name: "Should work",
			withNodegroupNames: []string{
				"ng-generic-1-21-1a-C29C1E3E88",
				"myname-1-20-3c-AAAAAAAAAA",
			},
			expectNodeGroups: []clusterConfigTemplateOptsNodeGroup{
				{
					Name:             "ng-generic-1-21-1a-C29C1E3E88",
					AvailabilityZone: "a",
				},
				{
					Name:             "myname-1-20-3c-AAAAAAAAAA",
					AvailabilityZone: "c",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := parseNodeGroups(tc.withNodegroupNames)

			assert.Equal(t, tc.expectNodeGroups, result)
		})
	}
}
