package eks_upgrade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func (e Eks) Version() (string, error) {
	args := []string{
		"get", "cluster", "-o", "json",
	}

	cmd := exec.Command("eksctl", args...) //nolint:gosec

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("starting eksctl command: %w. Details: %b", err, stderr)
	}

	var versionData VersionData

	err = json.Unmarshal(stdout, &versionData)
	if err != nil {
		return "", fmt.Errorf("unmarshalling version data: %w. Data: %s", err)
	}

}

func (e Eks) Upgrade() (string, error) {
	// create node group?
	// kubectl taint old nodes
}

func New() Eks {
	return Eks{}
}
