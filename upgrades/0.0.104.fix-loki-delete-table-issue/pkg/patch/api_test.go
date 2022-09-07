package patch

import (
	"io"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestAddBucketVersioning(t *testing.T) {
	testCases := []struct {
		name         string
		withTemplate io.Reader
	}{
		{
			name:         "Should add expected fields",
			withTemplate: strings.NewReader(testTemplate),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := AddDeleteTablePermission(tc.withTemplate)
			assert.NoError(t, err)

			rawResult, err := io.ReadAll(result)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, rawResult)
		})
	}
}

const testTemplate = `AWSTemplateFormatVersion: 2010-09-09
Outputs:
  LokiDynamoDBServiceAccountPolicy:
    Export:
      Name:
        Fn::Sub: ${AWS::StackName}-LokiDynamoDBServiceAccountPolicy
    Value:
      Ref: LokiDynamoDBServiceAccountPolicy
Resources:
  LokiDynamoDBServiceAccountPolicy:
    Properties:
      Description: Service account policy for storing indexes in an DynamoDB table
      ManagedPolicyName: okctl-julius-LokiDynamoDBServiceAccountPolicy
      PolicyDocument:
        Statement:
        - Action:
          - dynamodb:BatchGetItem
          - dynamodb:BatchWriteItem
          - dynamodb:UntagResource
          - dynamodb:PutItem
          - dynamodb:DeleteItem
          - dynamodb:ListTagsOfResource
          - dynamodb:Query
          - dynamodb:UpdateItem
          - dynamodb:CreateTable
          - dynamodb:TagResource
          - dynamodb:DescribeTable
          - dynamodb:GetItem
          - dynamodb:UpdateTable
          Effect: Allow
          Resource:
          - arn:aws:dynamodb:eu-west-1:932360772598:table/okctl-julius-loki-index_*
        - Action:
          - dynamodb:ListTables
          Effect: Allow
          Resource:
          - arn:aws:dynamodb:eu-west-1:932360772598:table/*
        Version: 2012-10-17
    Type: AWS::IAM::ManagedPolicy`
