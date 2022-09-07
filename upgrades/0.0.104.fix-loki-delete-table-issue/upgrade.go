package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.104.fix-loki-delete-table-issue/pkg/cfn"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.104.fix-loki-delete-table-issue/pkg/lib/logger"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.104.fix-loki-delete-table-issue/pkg/lib/manifest"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.104.fix-loki-delete-table-issue/pkg/patch"
	"github.com/spf13/afero"
)

func doUpgrade(ctx context.Context, log logger.Logger, fs *afero.Afero, dryRun bool) error {
	log.Debug("Acquiring cluster manifest information")

	clusterManifest, err := manifest.Cluster(fs)
	if err != nil {
		return fmt.Errorf("acquiring manifest: %w", err)
	}

	stackName := lokiDynamoDBPolicyStackName(clusterManifest.Metadata.Name)

	log.Debug("Fetching current Loki DynamoDB policy stack template")

	template, err := cfn.FetchStackTemplate(ctx, stackName)
	if err != nil {
		return fmt.Errorf("fetching stack template: %w", err)
	}

	template, err = ensureNoPermission(template, dynamoDBDeleteTablePermission)
	if err != nil {
		if errors.Is(err, errHasPermission) {
			log.Debug("Found relevant permission, ignoring upgrade")

			return nil
		}

		return fmt.Errorf("checking for existing permission: %w", err)
	}

	log.Debug("Adding dynamodb:DeleteTable permission")

	updatedTemplate, err := patch.AddDeleteTablePermission(template)
	if err != nil {
		return fmt.Errorf("updating stack template: %w", err)
	}

	log.Debug("Uploading patched template")

	if !dryRun {
		err = cfn.UpdateStackTemplate(ctx, stackName, updatedTemplate)
		if err != nil {
			return fmt.Errorf("uploading stack template: %w", err)
		}
	}

	log.Debug("Upgrade complete")

	return nil
}

func lokiDynamoDBPolicyStackName(clusterName string) string {
	return fmt.Sprintf("okctl-dynamodbpolicy-%s-loki", clusterName)
}

func ensureNoPermission(template io.Reader, permission string) (io.Reader, error) {
	raw, err := io.ReadAll(template)
	if err != nil {
		return nil, fmt.Errorf("buffering: %w", err)
	}

	if strings.Contains(string(raw), permission) {
		return nil, errHasPermission
	}

	return bytes.NewReader(raw), nil
}

var errHasPermission = errors.New("has permission")

const dynamoDBDeleteTablePermission = "dynamodb:DeleteTable"
