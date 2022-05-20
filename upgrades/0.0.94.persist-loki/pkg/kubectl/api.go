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

func GetLokiConfig(fs *afero.Afero) (io.Reader, error) {
	stdout, err := runCommand(fs,
		"--namespace", defaultMonitoringNamespace,
		"--output", "json",
		"get", secretResourceKind,
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

	lokiConfigAsString, ok := secret.Data[defaultLokiConfigSecretKey].(string)
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
		"patch", secretResourceKind, "loki",
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
		"get", podResourceKind, defaultLokiPodName,
	)
	if err != nil {
		if isErrNotFound(err) {
			return false, nil
		}

		return false, fmt.Errorf("running command: %w", err)
	}

	return true, nil
}

func RestartLoki(fs *afero.Afero) error {
	_, err := runCommand(fs,
		"--namespace", defaultMonitoringNamespace,
		"delete", podResourceKind, defaultLokiPodName,
	)
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}

func HasVolumeClaim(fs *afero.Afero, claimName string) (bool, error) {
	result, err := runCommand(fs,
		"--namespace", defaultMonitoringNamespace,
		"--output", "json",
		"get", persistentVolumeClaimResourceKind,
	)
	if err != nil {
		return false, fmt.Errorf("acquiring volume claims: %w", err)
	}

	resultAsBytes, err := io.ReadAll(result)
	if err != nil {
		return false, fmt.Errorf("buffering: %w", err)
	}

	var itemIterator iterator

	err = json.Unmarshal(resultAsBytes, &itemIterator)
	if err != nil {
		return false, fmt.Errorf("unmarshalling: %w", err)
	}

	for _, item := range itemIterator.Items {
		if item.Metadata.Name == claimName {
			return true, nil
		}
	}

	return false, nil
}

func GetVolumeClaimAZ(fs *afero.Afero, claimName string) (string, error) {
	volumeID, err := getVolumeID(fs, claimName)
	if err != nil {
		return "", fmt.Errorf("acquiring volume ID: %w", err)
	}

	zone, err := getVolumeZone(fs, volumeID)
	if err != nil {
		return "", fmt.Errorf("acquiring volume zone: %w", err)
	}

	return zone, nil
}

func AddNodeSelector(fs *afero.Afero, statefulsetName string, key string, value string) error {
	patch := jsp.New().Add(jsp.Operation{
		Type:  jsp.OperationTypeAdd,
		Path:  fmt.Sprintf("/spec/template/spec/nodeSelector/%s", key),
		Value: value,
	})

	rawPatch, err := patch.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}

	_, err = runCommand(fs,
		"--namespace", defaultMonitoringNamespace,
		"--type=json",
		"patch", statefulSetResourceKind, statefulsetName,
		"--patch", string(rawPatch),
	)
	if err != nil {
		return fmt.Errorf("patching: %w", err)
	}

	return nil
}
