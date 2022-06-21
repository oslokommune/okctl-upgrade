package upgrade

const (
	argoCDApplicationAPIVersion = "argoproj.io/v1alpha1"
	argoCDApplicationKind       = "Application"
)

type argoCDApplication struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

func (receiver argoCDApplication) Valid() bool {
	if receiver.ApiVersion != argoCDApplicationAPIVersion {
		return false
	}

	if receiver.Kind != argoCDApplicationKind {
		return false
	}

	return true
}
