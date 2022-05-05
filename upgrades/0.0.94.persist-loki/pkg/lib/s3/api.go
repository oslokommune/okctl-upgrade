package s3

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"time"
)

func CreateBucket(ctx context.Context, clusterName string) (string, error) {
	bucketName := fmt.Sprintf("okctl-%s-loki", clusterName)

	bucketTemplate, err := generateTemplate(bucketName)
	if err != nil {
		return "", fmt.Errorf("generating template: %w", err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("preparing config: %w", err)
	}

	client := cloudformation.NewFromConfig(cfg)
	waiter := cloudformation.NewStackCreateCompleteWaiter(client)
	stackName := fmt.Sprintf("okctl-s3bucket-%s-%s", clusterName, bucketName)

	_, err = client.CreateStack(ctx, &cloudformation.CreateStackInput{
		StackName:        aws.String(stackName),
		Tags:             generateTags(clusterName),
		TemplateBody:     aws.String(bucketTemplate),
		TimeoutInMinutes: aws.Int32(defaultStackTimeoutMinutes),
	})
	if err != nil {
		return "", fmt.Errorf("creating stack: %w", err)
	}

	out, err := waiter.WaitForOutput(
		ctx,
		&cloudformation.DescribeStacksInput{StackName: aws.String(stackName)},
		time.Minute*defaultStackTimeoutMinutes,
	)

	arn, err := getOutput(out, "S3Bucket", "BucketARN")
	if err != nil {
		return "", fmt.Errorf("getting bucket ARN: %w", err)
	}

	return arn, nil
}
