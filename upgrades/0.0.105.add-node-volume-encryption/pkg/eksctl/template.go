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
	AvailabilityZone string // 1a // 1b // 1c
}

func parseNodeGroups(nodegroupNames []string) []clusterConfigTemplateOptsNodeGroup {
	nodeGroups := make([]clusterConfigTemplateOptsNodeGroup, len(nodegroupNames))

	for index, item := range nodegroupNames {
		parts := strings.Split(item, "-")
		lastPart := parts[len(parts)-1]

		nodeGroups[index] = clusterConfigTemplateOptsNodeGroup{
			Name:             item,
			AvailabilityZone: string(lastPart[len(lastPart)-1]),
		}
	}

	return nodeGroups
}
