package upgrade

import (
	"fmt"

	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

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

func isPathAnArgoCDApplication(fs *afero.Afero, absoluteApplicationPath string) (bool, error) {
	stat, err := fs.Stat(absoluteApplicationPath)
	if err != nil {
		return false, fmt.Errorf("stating file: %w", err)
	}

	if stat.IsDir() {
		return false, nil
	}

	rawFile, err := fs.ReadFile(absoluteApplicationPath)
	if err != nil {
		return false, fmt.Errorf("reading file: %w", err)
	}

	var potentialApp argoCDApplication

	err = yaml.Unmarshal(rawFile, &potentialApp)
	if err != nil {
		return false, fmt.Errorf("unmarshalling: %w", err)
	}

	return potentialApp.Valid(), nil
}
