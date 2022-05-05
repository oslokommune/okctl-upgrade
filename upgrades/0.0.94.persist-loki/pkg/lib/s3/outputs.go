package s3

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

func getOutput(result *cloudformation.DescribeStacksOutput, _ string, outputName string) (string, error) {
	if len(result.Stacks) != 1 {
		return "", errors.New("unexpected amount of stacks")
	}

	for _, output := range result.Stacks[0].Outputs {
		if *output.OutputKey == outputName {
			return *output.OutputValue, nil
		}
	}

	return "", errors.New("output not found")
}
