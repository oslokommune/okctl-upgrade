package kubectl

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	jsp "github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/jsonpatch"
	"github.com/spf13/afero"
)

type Secret struct {
	Data map[string]interface{}
}

func GetLokiConfig(fs *afero.Afero) (io.Reader, error) {
	stdout, err := runCommand(fs,
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

	decodedLokiConfig, err := base64.StdEncoding.DecodeString(lokiConfigAsString)
	if err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}

	return bytes.NewReader(decodedLokiConfig), nil
}

func UpdateLokiConfig(fs *afero.Afero, config io.Reader) error {
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

	_, err = runCommand(fs,
		"--namespace", defaultMonitoringNamespace,
		"--type=json",
		"patch", "secret", "loki",
		"--patch", string(patchAsBytes),
	)
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}

func HasLoki(fs *afero.Afero) (bool, error) {
	_, err := runCommand(fs,
		"--namespace", defaultMonitoringNamespace,
		"--output", "json",
		"get", "pod", "loki-0",
	)
	if err != nil {
		if isErrNotFound(err) {
			return false, nil
		}

		return false, fmt.Errorf("running command: %w", err)
	}

	return true, nil
}
