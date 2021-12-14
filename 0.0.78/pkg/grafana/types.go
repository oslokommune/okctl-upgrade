package grafana

const jsonPatchOperationReplace = "replace"

type Patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

type Receipt string

type Receipts struct {
	receipts []Receipt
}
