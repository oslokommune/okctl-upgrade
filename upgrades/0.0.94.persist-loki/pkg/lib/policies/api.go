package policies

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/cfn"
)

const (
	defaultStackTimeoutMinutes     = 5
	defaultLogicalBucketPolicyName = "LokiS3ServiceAccountPolicy"
	defaultBucketPolicyOutputName  = "LokiS3ServiceAccountPolicy"
)

func CreateS3BucketPolicy(ctx context.Context, clusterName string, bucketARN string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("preparing config: %w", err)
	}

	client := cloudformation.NewFromConfig(cfg)
	stackName := fmt.Sprintf("okctl-s3bucketpolicy-%s-loki", clusterName)

	err = createBucketPolicyStack(ctx, client, clusterName, stackName, bucketARN)
	if err != nil {
		return "", fmt.Errorf("creating stack: %w", err)
	}

	out, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)})
	if err != nil {
		return "", fmt.Errorf("describing stack: %w", err)
	}

	arn, err := cfn.GetOutput(out, defaultLogicalBucketPolicyName, defaultBucketPolicyOutputName)
	if err != nil {
		return "", fmt.Errorf("getting ARN: %w", err)
	}

	return arn, nil
}
