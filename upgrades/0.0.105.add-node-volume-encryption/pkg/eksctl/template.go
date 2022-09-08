package eksctl

import (
	_ "embed"
	"strings"
)

//go:embed clusterconfig.yaml
var clusterConfigTemplate string

type clusterConfigTemplateOpts struct {
	ClusterName string
	AccountID   string
	Region      string
	NodeGroups  []clusterConfigTemplateOptsNodeGroup
}

type clusterConfigTemplateOptsNodeGroup struct {
	Name             string
	AvailabilityZone string
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
