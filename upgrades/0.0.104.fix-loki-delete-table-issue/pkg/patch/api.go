package patch

import (
	"bytes"
	"fmt"
	"io"
	"path"

	jsonpatch "github.com/evanphx/json-patch/v5"
	jsp "github.com/oslokommune/okctl-upgrade/upgrades/0.0.96.remote-state-versioning/pkg/lib/jsonpatch"

	"sigs.k8s.io/yaml"
)

var statementActionPath = path.Join(
	"/",
	"Resources",
	"LokiDynamoDBServiceAccountPolicy",
	"Properties",
	"PolicyDocument",
	"Statement",
	"0",
	"Action",
	"-",
)

const dynamoDBDeleteTablePermission = "dynamodb:DeleteTable"

func AddDeleteTablePermission(template io.Reader) (io.Reader, error) {
	rawTemplate, err := io.ReadAll(template)
	if err != nil {
		return nil, fmt.Errorf("buffering: %w", err)
	}

	rawTemplateAsJSON, err := yaml.YAMLToJSON(rawTemplate)
	if err != nil {
		return nil, fmt.Errorf("converting to JSON: %w", err)
	}

	patch := jsp.New().Add(
		jsp.Operation{
			Type:  jsp.OperationTypeAdd,
			Path:  statementActionPath,
			Value: dynamoDBDeleteTablePermission,
		},
	)

	rawPatch, err := patch.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshalling: %w", err)
	}

	decodedPatch, err := jsonpatch.DecodePatch(rawPatch)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	updatedTemplate, err := decodedPatch.Apply(rawTemplateAsJSON)
	if err != nil {
		return nil, fmt.Errorf("applying patch: %w", err)
	}

	updatedTemplateAsYAML, err := yaml.JSONToYAML(updatedTemplate)
	if err != nil {
		return nil, fmt.Errorf("converting to JSON: %w", err)
	}

	return bytes.NewReader(updatedTemplateAsYAML), nil
}
