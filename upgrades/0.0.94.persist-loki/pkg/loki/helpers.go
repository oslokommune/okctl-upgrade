package loki

import (
	"bytes"
	"fmt"
	"io"

	"sigs.k8s.io/yaml"
)

func asYAML(r io.Reader) (io.Reader, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("buffering: %w", err)
	}

	rawYAML, err := yaml.JSONToYAML(raw)
	if err != nil {
		return nil, fmt.Errorf("converting to YAML: %w", err)
	}

	return bytes.NewReader(rawYAML), nil
}
