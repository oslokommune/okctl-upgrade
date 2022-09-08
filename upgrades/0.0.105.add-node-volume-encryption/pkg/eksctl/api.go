package eksctl

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/manifest/apis/okctl.io/v1alpha1"
)

func GetClusterConfig(cluster v1alpha1.Cluster, nodegroupNames []string) (io.Reader, error) {
	t, err := template.New("clusterconfig").Parse(clusterConfigTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, clusterConfigTemplateOpts{
		ClusterName: cluster.Metadata.Name,
		AccountID:   cluster.Metadata.AccountID,
		Region:      cluster.Metadata.Region,
		NodeGroups:  parseNodeGroups(nodegroupNames),
	})
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &buf, nil
}

func parseNodeGroups(nodegroupNames []string) []clusterConfigTemplateOptsNodeGroup {
	nodeGroups := make([]clusterConfigTemplateOptsNodeGroup, len(nodegroupNames))

	for index, item := range nodegroupNames {
		parts := strings.Split(item, "-")

		nodeGroups[index] = clusterConfigTemplateOptsNodeGroup{
			Name:             item,
			AvailabilityZone: parts[len(parts)-1],
		}
	}

	return nodeGroups
}
