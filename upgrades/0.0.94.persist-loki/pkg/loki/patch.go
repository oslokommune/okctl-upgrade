package loki

import (
	"bytes"
	"fmt"
	"io"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	jsp "github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/jsonpatch"
)

func generateLokiPersistencePatch(region string, clusterName string, bucketName string, from time.Time) (io.Reader, error) {
	patch := jsp.New()

	patch.Add(
		jsp.Operation{
			Type:  jsp.OperationTypeAdd,
			Path:  "/schema_config/configs/-",
			Value: createS3SchemaConfig(clusterName, from),
		},
		jsp.Operation{
			Type:  jsp.OperationTypeAdd,
			Path:  "/storage_config/aws",
			Value: createAWSStorageConfig(region, bucketName),
		},
		jsp.Operation{
			Type:  jsp.OperationTypeReplace,
			Path:  "/table_manager",
			Value: createTableManagerIndexTablesProvisioning(),
		},
	)

	rawPatch, err := patch.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshalling patch: %w", err)
	}

	return bytes.NewReader(rawPatch), nil
}

func patchConfig(original io.Reader, patch io.Reader) (io.Reader, error) {
	rawPatch, err := io.ReadAll(patch)
	if err != nil {
		return nil, fmt.Errorf("buffering patch: %w", err)
	}

	decodedPatch, err := jsonpatch.DecodePatch(rawPatch)
	if err != nil {
		return nil, fmt.Errorf("decoding patch: %w", err)
	}

	configAsBytes, err := io.ReadAll(original)
	if err != nil {
		return nil, fmt.Errorf("buffering config: %w", err)
	}

	rawUpdatedConfig, err := decodedPatch.Apply(configAsBytes)
	if err != nil {
		return nil, fmt.Errorf("applying patch: %w", err)
	}

	return bytes.NewReader(rawUpdatedConfig), nil
}

func createS3SchemaConfig(clusterName string, from time.Time) SchemaConfig {
	return SchemaConfig{
		From:        from.Format("2006-01-02"),
		Store:       "aws",
		ObjectStore: "s3",
		Schema:      "v11",
		Index: SchemaConfigIndex{
			Prefix: fmt.Sprintf("okctl-%s-loki-index_", clusterName),
			Period: "336h",
		},
	}
}

func createAWSStorageConfig(region string, bucketName string) StorageConfig {
	return StorageConfig{
		S3:          fmt.Sprintf("s3://%s", region),
		BucketNames: bucketName,
		DynamoDB: map[string]string{
			"dynamodb_url": fmt.Sprintf("dynamodb://%s", region),
		},
	}
}

func createTableManagerIndexTablesProvisioning() TableManager {
	return TableManager{
		RetentionDeletesEnabled: true,
		RetentionPeriod:         "1344h",
		IndexTablesProvisioning: TableManagerIndexTablesProvisioning{
			ProvisionedWriteThroughput: 1,
			ProvisionedReadThroughput:  1,
			InactiveWriteThroughput:    1,
			InactiveReadThroughput:     1,
		},
	}
}
