package main

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/s3"
)

func upgrade(ctx context.Context, _ cmdflags.Flags) error {
	arn, err := s3.CreateBucket(ctx, "upgrade-test")
	if err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}

	fmt.Println(arn)

	//s3PolicyARN, err := something.CreateS3BucketPolicy(arn)
	//if err != nil {
	//	return fmt.Errorf("creating bucket policy: %w", err)
	//}

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
