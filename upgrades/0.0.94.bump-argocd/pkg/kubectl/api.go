package kubectl

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.bump-argocd/pkg/jsonpatch"
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
	serverContainerIndex, err := acquireContainerIndex(selector)
	if err != nil {
		return fmt.Errorf("acquiring container index: %w", err)
	}

	patch := jsonpatch.New().Add(jsonpatch.Operation{
		Type:  jsonpatch.OperationTypeReplace,
		Path:  fmt.Sprintf("/spec/template/spec/containers/%d/image", serverContainerIndex),
		Value: fmt.Sprintf("quay.io/argoproj/argocd:v%s", newVersion.String()),
	})

	rawPatch, err := patch.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}

	_, err = runCommand(
		"--namespace", selector.Namespace,
		"patch", selector.Kind, selector.Name,
		"--type", "json",
		"--patch", string(rawPatch),
	)
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}

func acquireContainerIndex(selector Selector) (int, error) {
	result, err := runCommand(
		"--namespace", selector.Namespace,
		"--output", "json",
		"get", selector.Kind, selector.Name,
	)
	if err != nil {
		return -1, fmt.Errorf("running command: %w", err)
	}

	raw, err := io.ReadAll(result)
	if err != nil {
		return -1, fmt.Errorf("buffering: %w", err)
	}

	var resource resources.Deployment

	err = json.Unmarshal(raw, &resource)
	if err != nil {
		return -1, fmt.Errorf("unmarshalling: %w", err)
	}

	for index, container := range resource.Spec.Template.Spec.Containers {
		if container.Name == selector.ContainerName {
			return index, nil
		}
	}

	return -1, fmt.Errorf("no container with name %s found", selector.ContainerName)
}
