package kubectl

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	jsp "github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/jsonpatch"
	"github.com/spf13/afero"
)

type Secret struct {
	Data map[string]interface{}
}

func GetLokiConfig(fs *afero.Afero, clusterName string) (io.Reader, error) {
	kubeconfigPath, err := acquireKubeconfigPath(clusterName)
	if err != nil {
		return nil, fmt.Errorf("acquiring kubeconfig path: %w", err)
	}

	stdout, err := runCommand(fs, kubeconfigPath,
		"--namespace", defaultMonitoringNamespace,
		"--output", "json",
		"get", "secret",
		"loki",
	)
	if err != nil {
		return nil, fmt.Errorf("running command: %w", err)
	}

	rawSecret, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("buffering secret: %w", err)
	}

	var secret Secret

	err = json.Unmarshal(rawSecret, &secret)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling secret: %w", err)
	}

	lokiConfigAsString, ok := secret.Data["loki.yaml"].(string)
	if !ok {
		return nil, fmt.Errorf("converting config to string")
	}

	return strings.NewReader(lokiConfigAsString), nil
}

func UpdateLokiConfig(fs *afero.Afero, clusterName string, config io.Reader) error {
	kubeconfigPath, err := acquireKubeconfigPath(clusterName)
	if err != nil {
		return fmt.Errorf("acquiring kubeconfig path: %w", err)
	}

	rawConfig, err := io.ReadAll(config)
	if err != nil {
		return fmt.Errorf("buffering config: %w", err)
	}

	p := jsp.New()
	p.Add(
		jsp.Operation{
			Type:  jsp.OperationTypeReplace,
			Path:  "/data/loki.yaml",
			Value: rawConfig,
		},
	)

	patchAsBytes, err := p.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshalling patch: %w", err)
	}

	_, err = runCommand(fs, kubeconfigPath,
		"--namespace", defaultMonitoringNamespace,
		"--type='json'",
		"patch", "secret", "loki",
		fmt.Sprintf("--patch='%s'", string(patchAsBytes)),
	)
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}
