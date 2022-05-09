package loki

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"testing"

	"sigs.k8s.io/yaml"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/pre-persistence-values.yaml
var prePersistenceValues string

func TestPatch(t *testing.T) {
	testCases := []struct {
		name         string
		withOriginal string
	}{
		{
			name:         "Should properly patch values",
			withOriginal: prePersistenceValues,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			originalAsJSON, err := yaml.YAMLToJSON([]byte(tc.withOriginal))
			assert.NoError(t, err)

			result, err := patchValues(
				bytes.NewReader(originalAsJSON),
				"eu-test-1",
				"mock-cluster",
				"mock-bucket",
			)
			assert.NoError(t, err)

			rawResult, err := io.ReadAll(result)
			assert.NoError(t, err)

			prettyResult := bytes.Buffer{}

			err = json.Indent(&prettyResult, rawResult, "", "  ")
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, prettyResult.Bytes())
		})
	}
}

func TestPatchGeneration(t *testing.T) {
	testCases := []struct {
		name            string
		withRegion      string
		withClusterName string
		withBucketName  string
	}{
		{
			name:            "Should produce expected JSON patch data",
			withRegion:      "eu-test-1",
			withClusterName: "mock-cluster",
			withBucketName:  "mock-bucket",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			patch, err := generateLokiPersistencePatch(tc.withRegion, tc.withClusterName, tc.withBucketName)
			assert.NoError(t, err)

			rawPatch, err := io.ReadAll(patch)
			assert.NoError(t, err)

			prettyPrinted := bytes.Buffer{}

			err = json.Indent(&prettyPrinted, rawPatch, "", "  ")
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, prettyPrinted.Bytes())
		})
	}
}
