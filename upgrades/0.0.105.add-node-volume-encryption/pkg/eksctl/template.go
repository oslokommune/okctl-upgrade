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
	Name             string // ng-generic-1-21-1a-C29C1E3E88
	AvailabilityZone string // a // b // c
}

func parseNodeGroups(nodegroupNames []string) []clusterConfigTemplateOptsNodeGroup {
	nodeGroups := make([]clusterConfigTemplateOptsNodeGroup, len(nodegroupNames))

	for index, item := range nodegroupNames {
		parts := reverse(strings.Split(item, "-"))
		azContainer := parts[1]

		nodeGroups[index] = clusterConfigTemplateOptsNodeGroup{
			Name:             item,
			AvailabilityZone: string(azContainer[len(azContainer)-1]),
		}
	}

	return nodeGroups
}

func reverse(l []string) []string {
	reversedList := make([]string, len(l))
	reversedIndex := len(l) - 1

	for _, item := range l {
		reversedList[reversedIndex] = item

		reversedIndex--
	}

	return reversedList
}
