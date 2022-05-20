package loki

type SchemaConfig struct {
	From        string            `json:"from"`
	Store       string            `json:"store"`
	ObjectStore string            `json:"object_store"`
	Schema      string            `json:"schema"`
	Index       SchemaConfigIndex `json:"index"`
}

type SchemaConfigIndex struct {
	Prefix string `json:"prefix"`
	Period string `json:"period"`
}

type StorageConfig struct {
	S3          string            `json:"s3"`
	BucketNames string            `json:"bucketnames"`
	DynamoDB    map[string]string `json:"dynamodb"`
}

type TableManagerIndexTablesProvisioning struct {
	EnableOnDemandThroughputMode         bool `json:"enable_ondemand_throughput_mode"`
	EnableInactiveThroughputOnDemandMode bool `json:"enable_inactive_throughput_on_demand_mode"`
}

type TableManager struct {
	RetentionDeletesEnabled bool                                `json:"retention_deletes_enabled"`
	RetentionPeriod         string                              `json:"retention_period"`
	IndexTablesProvisioning TableManagerIndexTablesProvisioning `json:"index_tables_provisioning"`
}

type ServiceAccountConfig struct {
	Create      bool              `json:"create"`
	Name        string            `json:"name"`
	Annotations map[string]string `json:"annotations"`
}
