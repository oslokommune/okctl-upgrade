package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func createBucket(cfg aws.Config, clusterName string) (string, error) {
	cfn := cloudformation.NewFromConfig(cfg)

	cfnTemplate, err := generateBucketCFN(clusterName)
	if err != nil {
		return "", fmt.Errorf("generating CloudFormation template: %w", err)
	}

	stackName := generateBucketStackName(clusterName)

	_, err = cfn.CreateStack(context.Background(), &cloudformation.CreateStackInput{
		StackName: aws.String(stackName),
		Tags: []types.Tag{
			{
				Key:   aws.String("alpha.okctl.io/cluster-name"),
				Value: aws.String(clusterName),
			},
			{
				Key:   aws.String("alpha.okctl.io/managed"),
				Value: aws.String("true"),
			},
		},
		TemplateBody:     aws.String(cfnTemplate),
		TimeoutInMinutes: aws.Int32(5),
	})
	if err != nil {
		var (
			aerr    *types.AlreadyExistsException
			authErr smithy.APIError
		)

		if errors.As(err, &aerr) {
			return "", ErrAlreadyExists
		}

		if errors.As(err, &authErr) {
			return "", ErrNotAuthenticated
		}

		return "", fmt.Errorf("creating bucket: %w", err)
	}

	bucketID, err := acquireID(cfn, stackName)
	if err != nil {
		return "", fmt.Errorf("acquiring ARN: %w", err)
	}

	return bucketID, nil
}

func acquireID(cfn *cloudformation.Client, stackName string) (string, error) {
	for i := 1; i < 4; i++ {
		duration := time.Duration(15*i) * time.Second

		time.Sleep(duration)

		result, err := cfn.DescribeStacks(context.Background(), &cloudformation.DescribeStacksInput{
			StackName: aws.String(stackName),
		})
		if err != nil {
			return "", fmt.Errorf("describing stack: %w", err)
		}

		if result.Stacks[0].StackStatus == types.StackStatusCreateComplete {
			outputs := outputsAsMap(result.Stacks[0].Outputs)

			return outputs["S3Bucket"], nil
		}

		log.Printf("CloudFormation stack not ready after %d seconds, retrying", duration)
	}

	return "", errors.New("waiting for stack to complete")
}

func generateBucketStackName(clusterName string) string {
	return fmt.Sprintf("okctl-s3bucket-%s-okctl-%s-meta", clusterName, clusterName)
}

func generateBucketCFN(clusterName string) (string, error) {
	result := bytes.Buffer{}
	bucketTemplate, _ := template.New("bucketTemplateRaw").Parse(bucketTemplateRaw)

	err := bucketTemplate.Execute(&result, bucketTemplateOpts{ClusterName: clusterName})
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return result.String(), nil
}

func outputsAsMap(outputs []types.Output) map[string]string {
	result := make(map[string]string)

	for _, output := range outputs {
		result[*output.OutputKey] = *output.OutputValue
	}

	return result
}

type bucketTemplateOpts struct {
	ClusterName string
}

const bucketTemplateRaw = `
AWSTemplateFormatVersion: 2010-09-09
Outputs:
  S3Bucket:
    Export:
      Name:
        Fn::Sub: ${AWS::StackName}-S3Bucket
    Value:
      Ref: S3Bucket
  BucketArn:
    Export:
      Name:
        Fn::Sub: ${AWS::StackName}-BucketArn
    Value:
      Fn::GetAtt: S3Bucket.Arn
Resources:
  S3Bucket:
    Properties:
      AccessControl: BucketOwnerFullControl
      BucketName: okctl-{{- .ClusterName -}}-meta
      VersioningConfiguration:
        Status: Enabled
    Type: AWS::S3::Bucket
`
