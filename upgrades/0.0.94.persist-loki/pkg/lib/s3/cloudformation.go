package s3

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"text/template"
	"time"
)

func createBucketStack(ctx context.Context, client *cloudformation.Client, clusterName, stackName, bucketName string) error {
	bucketTemplate, err := generateTemplate(bucketName)
	if err != nil {
		return fmt.Errorf("generating template: %w", err)
	}

	_, err = client.CreateStack(ctx, &cloudformation.CreateStackInput{
		StackName:        aws.String(stackName),
		Tags:             generateTags(clusterName),
		TemplateBody:     aws.String(bucketTemplate),
		TimeoutInMinutes: aws.Int32(defaultStackTimeoutMinutes),
	})
	if err != nil {
		var alreadyExists *types.AlreadyExistsException

		if errors.As(err, &alreadyExists) {
			return nil
		}

		return fmt.Errorf("creating stack: %w", err)
	}

	waiter := cloudformation.NewStackCreateCompleteWaiter(client)

	err = waiter.Wait(
		ctx,
		&cloudformation.DescribeStacksInput{StackName: aws.String(stackName)},
		time.Minute*defaultStackTimeoutMinutes,
	)
	if err != nil {
		return fmt.Errorf("waiting for stack: %w", err)
	}

	return nil
}

func generateTemplate(bucketName string) (string, error) {
	buf := bytes.Buffer{}

	t, err := template.New("bucket").Parse(rawBucketTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	err = t.Execute(&buf, struct {
		BucketName string
	}{
		BucketName: bucketName,
	})
	if err != nil {
		return "", fmt.Errorf("interpolating template: %w", err)
	}

	return buf.String(), nil
}

//go:embed bucket-template.yaml
var rawBucketTemplate string
