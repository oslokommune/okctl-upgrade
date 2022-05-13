package loki

import (
	"fmt"
	"time"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/context"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/kubectl"
)

func AddPersistence(ctx context.Context, region string, clusterName string, bucketName string) error {
	patch, err := generateLokiPersistencePatch(region, clusterName, bucketName, time.Now())
	if err != nil {
		return fmt.Errorf("generating persistence patch: %w", err)
	}

	original, err := kubectl.GetLokiConfig(ctx.Fs)
	if err != nil {
		return fmt.Errorf("acquiring config: %w", err)
	}

	originalAsJSON, err := asJSON(original)
	if err != nil {
		return fmt.Errorf("converting to JSON: %w", err)
	}

	updated, err := patchConfig(originalAsJSON, patch)
	if err != nil {
		return fmt.Errorf("patching config: %w", err)
	}

	updatedConfigAsYAML, err := asYAML(updated)
	if err != nil {
		return fmt.Errorf("converting to YAML: %w", err)
	}

	if ctx.Flags.DryRun {
		ctx.Logger.Debug("Patching config locally successful. Applying to cluster")

		return nil
	}

	err = kubectl.UpdateLokiConfig(ctx.Fs, updatedConfigAsYAML)
	if err != nil {
		return fmt.Errorf("updating config: %w", err)
	}

	err = kubectl.RestartLoki(ctx.Fs)
	if err != nil {
		return fmt.Errorf("restarting Loki: %w", err)
	}

	return nil
}
