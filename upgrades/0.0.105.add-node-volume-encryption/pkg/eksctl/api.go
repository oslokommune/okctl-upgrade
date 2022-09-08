package eksctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/manifest/apis/okctl.io/v1alpha1"
)

func GenerateClusterConfig(cluster v1alpha1.Cluster, nodegroupNames []string) (io.Reader, error) {
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

type getNodeGroupResult struct {
	Name string `json:"Name"`
}

func GetNodeGroupNames(clusterName string) ([]string, error) {
	result, err := runCommand("get", "nodegroup", "--cluster", clusterName, "-o", "json")
	if err != nil {
		return nil, fmt.Errorf("retrieving nodegroups: %w", err)
	}

	var parsedResult []getNodeGroupResult

	rawResult, err := io.ReadAll(result)
	if err != nil {
		return nil, fmt.Errorf("buffering result: %w", err)
	}

	err = json.Unmarshal(rawResult, &parsedResult)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	nodegroupNames := make([]string, len(parsedResult))

	for index, item := range parsedResult {
		nodegroupNames[index] = item.Name
	}

	return nodegroupNames, nil
}
