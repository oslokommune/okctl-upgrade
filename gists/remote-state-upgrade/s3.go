package main

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func uploadState(cfg aws.Config, bucketID string, reader io.Reader) error {
	client := s3.NewFromConfig(cfg)

	_, err := client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:               aws.String(bucketID),
		Key:                  aws.String("state.db"),
		Body:                 reader,
		ServerSideEncryption: "aws:kms",
	})
	if err != nil {
		return fmt.Errorf("putting object: %w", err)
	}

	return nil
}
