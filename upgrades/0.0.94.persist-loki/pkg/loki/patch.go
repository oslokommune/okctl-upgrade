package loki

import (
	"bytes"
	"fmt"
	"io"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	jsp "github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/jsonpatch"
)

func generateLokiPersistencePatch(region string, clusterName string, bucketName string) (io.Reader, error) {
	patch := jsp.New()

	patch.Add(
		jsp.Operation{
			Type:  jsp.OperationTypeAdd,
			Path:  "/config/schema_config/configs/-",
			Value: createS3SchemaConfig(clusterName),
		},
		jsp.Operation{
			Type:  jsp.OperationTypeAdd,
			Path:  "/config/storage_config/aws",
			Value: createAWSStorageConfig(region, bucketName),
		},
		jsp.Operation{
			Type:  jsp.OperationTypeReplace,
			Path:  "/config/table_manager",
			Value: createTableManagerIndexTablesProvisioning(),
		},
		jsp.Operation{
			Type:  jsp.OperationTypeReplace,
			Path:  "/serviceAccount",
			Value: createServiceAccountConfig(),
		},
	)

	rawPatch, err := patch.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshalling patch: %w", err)
	}

	return bytes.NewReader(rawPatch), nil
}

func patchValues(original io.Reader, region string, clusterName string, bucketName string) (io.Reader, error) {
	patch, err := generateLokiPersistencePatch(region, clusterName, bucketName)
	if err != nil {
		return nil, fmt.Errorf("generating patch: %w", err)
	}

	rawPatch, err := io.ReadAll(patch)
	if err != nil {
		return nil, fmt.Errorf("buffering patch: %w", err)
	}

	p, err := jsonpatch.DecodePatch(rawPatch)
	if err != nil {
		return nil, fmt.Errorf("decoding patch: %w", err)
	}

	rawValues, err := io.ReadAll(original)
	if err != nil {
		return nil, fmt.Errorf("buffering original: %w", err)
	}

	rawPatchedValues, err := p.Apply(rawValues)
	if err != nil {
		return nil, fmt.Errorf("patching values: %w", err)
	}

	return bytes.NewReader(rawPatchedValues), nil
}

func createS3SchemaConfig(clusterName string) SchemaConfig {
	return SchemaConfig{
		From:        time.Now().Format("2006-01-02"),
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

func createServiceAccountConfig() ServiceAccountConfig {
	return ServiceAccountConfig{
		Create:      false,
		Name:        "loki",
		Annotations: make(map[string]string),
	}
}
