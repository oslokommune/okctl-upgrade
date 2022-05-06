package main

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/eksctl"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/policies"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/s3"
	"github.com/spf13/afero"
)

func upgrade(ctx context.Context, _ cmdflags.Flags) error {
	fs := &afero.Afero{Fs: afero.NewOsFs()}
	clusterName := "upgrade-test"
	awsAccountID := ""
	awsRegion := "eu-west-1"

	bucketName := fmt.Sprintf("okctl-%s-loki", clusterName)

	arn, err := s3.CreateBucket(ctx, clusterName, bucketName)
	if err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}

	fmt.Printf("S3 ARN: %s\n", arn)

	s3PolicyARN, err := policies.CreateS3BucketPolicy(ctx, clusterName, arn)
	if err != nil {
		return fmt.Errorf("creating bucket policy: %w", err)
	}

	fmt.Printf("Bucket policy ARN: %s\n", s3PolicyARN)

	dynamoDBPolicyARN, err := policies.CreateDynamoDBPolicy(
		ctx,
		awsAccountID,
		awsRegion,
		clusterName,
	)
	if err != nil {
		return fmt.Errorf("creating Dynamo DB policy: %w", err)
	}

	fmt.Printf("DynamoDB policy ARN: %s\n", dynamoDBPolicyARN)

	err = eksctl.CreateServiceUser(fs, clusterName, "loki", []string{s3PolicyARN, dynamoDBPolicyARN})
	if err != nil {
		return fmt.Errorf("creating service user: %w", err)
	}

	//err = something.PatchLoki()
	//if err != nil {
	//	return fmt.Errorf("patching Loki: %w", err)
	//}

	return nil
}
