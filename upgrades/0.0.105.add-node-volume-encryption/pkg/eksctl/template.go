package eksctl

import _ "embed"

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
