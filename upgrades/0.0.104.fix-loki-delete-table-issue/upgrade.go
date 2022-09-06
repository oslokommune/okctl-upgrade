package main

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.104.fix-loki-delete-table-issue/pkg/cfn"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.104.fix-loki-delete-table-issue/pkg/lib/manifest"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.104.fix-loki-delete-table-issue/pkg/patch"
	"github.com/spf13/afero"
)

func doUpgrade(ctx context.Context, fs *afero.Afero) error {
	clusterManifest, err := manifest.Cluster(fs)
	if err != nil {
		return fmt.Errorf("acquiring manifest: %w", err)
	}

	stackName := lokiDynamoDBPolicyStackName(clusterManifest.Metadata.Name)

	template, err := cfn.FetchStackTemplate(ctx, stackName)
	if err != nil {
		return fmt.Errorf("fetching stack template: %w", err)
	}

	updatedTemplate, err := patch.AddDeleteTablePermission(template)
	if err != nil {
		return fmt.Errorf("updating stack template: %w", err)
	}

	err = cfn.UpdateStackTemplate(ctx, stackName, updatedTemplate)
	if err != nil {
		return fmt.Errorf("uploading stack template: %w", err)
	}

	return nil
}

func lokiDynamoDBPolicyStackName(clusterName string) string {
	return fmt.Sprintf("okctl-dynamodbpolicy-%s-loki", clusterName)
}
