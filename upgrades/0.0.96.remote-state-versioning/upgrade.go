package main

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.96.remote-state-versioning/pkg/cfn"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.96.remote-state-versioning/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.96.remote-state-versioning/pkg/lib/logging"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.96.remote-state-versioning/pkg/patch"
)

func upgrade(ctx context.Context, log logging.Logger, flags cmdflags.Flags) error {
	stackName := stackName("julius-one")

	log.Debugf("Fetching stack with name %s", stackName)

	template, err := cfn.FetchStackTemplate(ctx, stackName)
	if err != nil {
		return fmt.Errorf("fetching template: %w", err)
	}

	log.Debug("Found stack. Starting patch operation")

	patchedTemplate, err := patch.AddBucketVersioning(template)
	if err != nil {
		return fmt.Errorf("patching: %w", err)
	}

	log.Debug("Patching successful. Updating stack.")

	if !flags.DryRun {
		err = cfn.UpdateStackTemplate(ctx, stackName, patchedTemplate)
		if err != nil {
			return fmt.Errorf("updating template: %w", err)
		}
	}

	log.Debug("Update success.")

	return nil
}

func stackName(clusterName string) string {
	return fmt.Sprintf("okctl-s3bucket-%s-okctl-%s-meta", clusterName, clusterName)
}
