package main

import (
	"context"
	"fmt"
	"os"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/eksctl"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/loki"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/policies"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/s3"
	"github.com/spf13/afero"
)

func upgrade(ctx context.Context, _ cmdflags.Flags) error {
	fs := &afero.Afero{Fs: afero.NewOsFs()}
	clusterManifestPath := os.Getenv("OKCTL_CLUSTER_DECLARATION")

	rawManifest, err := fs.ReadFile(clusterManifestPath)
	if err != nil {
		return fmt.Errorf("reading manifest: %w", err)
	}

	clusterManifest, err := v1alpha1.Parse(rawManifest)
	if err != nil {
		return fmt.Errorf("parsing manifest: %w", err)
	}

	bucketName := fmt.Sprintf("okctl-%s-loki", clusterManifest.Metadata.Name)

	arn, err := s3.CreateBucket(ctx, clusterManifest.Metadata.Name, bucketName)
	if err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}

	fmt.Printf("S3 ARN: %s\n", arn)

	s3PolicyARN, err := policies.CreateS3BucketPolicy(ctx, clusterManifest.Metadata.Name, arn)
	if err != nil {
		return fmt.Errorf("creating bucket policy: %w", err)
	}

	fmt.Printf("Bucket policy ARN: %s\n", s3PolicyARN)

	dynamoDBPolicyARN, err := policies.CreateDynamoDBPolicy(
		ctx,
		clusterManifest.Metadata.AccountID,
		clusterManifest.Metadata.Region,
		clusterManifest.Metadata.Name,
	)
	if err != nil {
		return fmt.Errorf("creating Dynamo DB policy: %w", err)
	}

	fmt.Printf("DynamoDB policy ARN: %s\n", dynamoDBPolicyARN)

	err = eksctl.CreateServiceUser(
		fs,
		clusterManifest.Metadata.Name,
		"loki",
		[]string{s3PolicyARN, dynamoDBPolicyARN},
	)
	if err != nil {
		return fmt.Errorf("creating service user: %w", err)
	}

	err = loki.AddPersistence(
		fs,
		clusterManifest.Metadata.Region,
		clusterManifest.Metadata.Name,
		bucketName,
	)
	if err != nil {
		return fmt.Errorf("patching Loki: %w", err)
	}

	return nil
}
