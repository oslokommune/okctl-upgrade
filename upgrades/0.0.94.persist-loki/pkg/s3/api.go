package s3

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/cfn"
)

// CreateBucket knows how to make a bucket through CloudFormation
func CreateBucket(ctx context.Context, clusterName string, bucketName string) (string, error) {
	stackName := fmt.Sprintf("okctl-s3bucket-%s-%s", clusterName, bucketName)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("preparing config: %w", err)
	}

	client := cloudformation.NewFromConfig(cfg)

	err = createBucketStack(ctx, client, clusterName, stackName, bucketName)
	if err != nil {
		return "", fmt.Errorf("creating stack: %w", err)
	}

	out, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return "", fmt.Errorf("describing stack: %w", err)
	}

	arn, err := cfn.GetOutput(out, defaultLogicalBucketName, defaultBucketARNOutputName)
	if err != nil {
		return "", fmt.Errorf("getting bucket ARN: %w", err)
	}

	return arn, nil
}
