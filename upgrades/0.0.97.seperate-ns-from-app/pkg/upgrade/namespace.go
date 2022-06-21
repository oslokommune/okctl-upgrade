package upgrade

type namespaceMetadata struct {
	Name string
}

type namespaceManifest struct {
	ApiVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   namespaceMetadata `json:"metadata"`
}
