package kubectl

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.bump-argocd/pkg/kubectl/resources"
	"io"
	"strings"
)

type Selector struct {
	Namespace     string
	Kind          string
	Name          string
	ContainerName string
}

func GetImageVersion(selector Selector) (semver.Version, error) {
	result, err := runCommand(
		"--namespace", selector.Namespace,
		"--output", "json",
		"get", selector.Kind, selector.Name,
	)
	if err != nil {
		return semver.Version{}, fmt.Errorf("running command: %w", err)
	}

	raw, err := io.ReadAll(result)
	if err != nil {
		return semver.Version{}, fmt.Errorf("buffering: %w", err)
	}

	var resource resources.Deployment

	err = json.Unmarshal(raw, &resource)
	if err != nil {
		return semver.Version{}, fmt.Errorf("unmarshalling: %w", err)
	}

	var imageURI string

	for _, container := range resource.Spec.Template.Spec.Containers {
		if container.Name == selector.ContainerName {
			imageURI = container.Image
		}
	}

	if imageURI == "" {
		return semver.Version{}, fmt.Errorf("container name %s not found", selector.ContainerName)
	}

	tag := strings.SplitN(imageURI, ":", 2)[1]
	tag = strings.TrimPrefix(tag, "v")

	version, err := semver.NewVersion(tag)
	if err != nil {
		return semver.Version{}, fmt.Errorf("parsing tag: %w", err)
	}

	return *version, nil
}

func UpdateImageVersion(selector Selector, newVersion semver.Version) error {
	return nil
}
