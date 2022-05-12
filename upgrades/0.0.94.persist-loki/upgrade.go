package main

import (
	"fmt"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/context"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/eksctl"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/loki"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/policies"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/s3"
	"github.com/spf13/afero"
)

func upgrade(ctx context.Context, fs *afero.Afero, clusterManifest v1alpha1.Cluster) error {
	bucketName := fmt.Sprintf("okctl-%s-loki", clusterManifest.Metadata.Name)

	ctx.Logger.Debug("Creating S3 bucket")

	arn, err := s3.CreateBucket(ctx.Ctx, clusterManifest.Metadata.Name, bucketName)
	if err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}

	ctx.Logger.Debugf("Successfully created S3 bucket with ARN: %s\n", arn)

	ctx.Logger.Debug("Creating S3 bucket policy")

	s3PolicyARN, err := policies.CreateS3BucketPolicy(ctx.Ctx, clusterManifest.Metadata.Name, arn)
	if err != nil {
		return fmt.Errorf("creating bucket policy: %w", err)
	}

	ctx.Logger.Debugf("Successfully created bucket policy with ARN: %s\n", s3PolicyARN)

	ctx.Logger.Debug("Creating DynamoDB policy")

	dynamoDBPolicyARN, err := policies.CreateDynamoDBPolicy(
		ctx.Ctx,
		clusterManifest.Metadata.AccountID,
		clusterManifest.Metadata.Region,
		clusterManifest.Metadata.Name,
	)
	if err != nil {
		return fmt.Errorf("creating Dynamo DB policy: %w", err)
	}

	ctx.Logger.Debugf("Successfully created DynamoDB policy with ARN: %s\n", dynamoDBPolicyARN)

	ctx.Logger.Debug("Creating Kubernetes service user -> role mapping with relevant policies")

	err = eksctl.CreateServiceUser(
		fs,
		clusterManifest.Metadata.Name,
		"loki",
		[]string{s3PolicyARN, dynamoDBPolicyARN},
	)
	if err != nil {
		return fmt.Errorf("creating service user: %w", err)
	}

	ctx.Logger.Debug("Successfully created service user Loki")

	ctx.Logger.Debug("Patching Loki config to add new storage configuration")

	err = loki.AddPersistence(
		fs,
		clusterManifest.Metadata.Region,
		clusterManifest.Metadata.Name,
		bucketName,
	)
	if err != nil {
		return fmt.Errorf("patching Loki: %w", err)
	}

	ctx.Logger.Debug("Successfully patched Loki config")

	return nil
}
