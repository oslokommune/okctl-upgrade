package loki

import (
	"fmt"
	"time"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/kubectl"
	"github.com/spf13/afero"
)

func AddPersistence(fs *afero.Afero, kubeconfigPath string, region string, clusterName string, bucketName string) error {
	patch, err := generateLokiPersistencePatch(region, clusterName, bucketName, time.Now())
	if err != nil {
		return fmt.Errorf("generating persistence patch: %w", err)
	}

	originalConfig, err := kubectl.GetLokiConfig(fs, kubeconfigPath)
	if err != nil {
		return fmt.Errorf("acquiring config: %w", err)
	}

	updatedConfig, err := patchConfig(originalConfig, patch)
	if err != nil {
		return fmt.Errorf("patching config: %w", err)
	}

	updatedConfigAsYAML, err := asYAML(updatedConfig)
	if err != nil {
		return fmt.Errorf("converting to YAML: %w", err)
	}

	err = kubectl.UpdateLokiConfig(fs, kubeconfigPath, updatedConfigAsYAML)
	if err != nil {
		return fmt.Errorf("updating config: %w", err)
	}

	return nil
}
