package policies

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/cfn"
	"text/template"
	"time"
)

func createBucketPolicyStack(ctx context.Context, client *cloudformation.Client, clusterName string, stackName string, bucketARN string) error {
	policyTemplate, err := generateBucketPolicyTemplate(clusterName, bucketARN)
	if err != nil {
		return fmt.Errorf("generating template: %w", err)
	}

	_, err = client.CreateStack(ctx, &cloudformation.CreateStackInput{
		StackName:        aws.String(stackName),
		Tags:             cfn.GenerateTags(clusterName),
		TemplateBody:     aws.String(policyTemplate),
		TimeoutInMinutes: aws.Int32(defaultStackTimeoutMinutes),
		Capabilities:     []types.Capability{types.CapabilityCapabilityNamedIam},
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

func generateBucketPolicyTemplate(clusterName string, bucketARN string) (string, error) {
	buf := bytes.Buffer{}

	t, err := template.New("bucket-policy").Parse(rawBucketPolicyTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	err = t.Execute(&buf, struct {
		ClusterName string
		BucketARN   string
	}{
		ClusterName: clusterName,
		BucketARN:   bucketARN,
	})
	if err != nil {
		return "", fmt.Errorf("interpolating template: %w", err)
	}

	return buf.String(), nil
}

//go:embed bucket-policy.yaml
var rawBucketPolicyTemplate string
