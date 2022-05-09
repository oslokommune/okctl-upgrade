package loki

import (
	"bytes"
	_ "embed"
	"io"
	"testing"

	"sigs.k8s.io/yaml"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

//go:embed pre-persistence-values.yaml
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
				"mock-cluster",
			)
			assert.NoError(t, err)

			rawResult, err := io.ReadAll(result)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, rawResult)
		})
	}
}
