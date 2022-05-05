package main

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/policies"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/s3"
)

func upgrade(ctx context.Context, _ cmdflags.Flags) error {
	bucketName := fmt.Sprintf("okctl-%s-loki", "upgrade-test")
	clusterName := "upgrade-test"

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

	//dynamoDBPolicyARN, err := something.CreateDynamoDBPolicy()
	//if err != nil {
	//	return fmt.Errorf("creating Dynamo DB policy: %w", err)
	//}

	//err = something.CreateServiceUser("loki", s3PolicyARN, dynamoDBPolicyARN)
	//if err != nil {
	//	return fmt.Errorf("creating service user: %w", err)
	//}

	//err = something.PatchLoki()
	//if err != nil {
	//	return fmt.Errorf("patching Loki: %w", err)
	//}

	return nil
}
