package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func generateTags(clusterName string) []types.Tag {
	return []types.Tag{
		{
			Key:   aws.String("alpha.okctl.io/cluster-name"),
			Value: aws.String(clusterName),
		},
		{
			Key:   aws.String("alpha.okctl.io/managed"),
			Value: aws.String("true"),
		},
		{
			Key:   aws.String("alpha.okctl.io/okctl-commit"),
			Value: aws.String("unknown"),
		},
		{
			Key:   aws.String("alpha.okctl.io/okctl-version"),
			Value: aws.String("0.0.94"),
		},
	}
}
