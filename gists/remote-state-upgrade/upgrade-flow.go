package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
)

func doUpgrade(clusterManifest Cluster) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(fmt.Errorf("generating AWS config: %w", err))
	}

	log.Println("Creating bucket")

	bucketID, err := createBucket(cfg, clusterManifest.Metadata.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrAlreadyExists):
			log.Fatal(t(BucketAlreadyExists))
		case errors.Is(err, ErrNotAuthenticated):
			log.Fatal(t(NotAuthenticated))
		default:
			log.Fatal(fmt.Errorf("creating bucket: %w", err))
		}
	}

	log.Println("Uploading state.db")

	stateDbContents, err := loadStateDB(clusterManifest)
	if err != nil {
		log.Fatal(fmt.Errorf("loading state database: %w", err))
	}

	err = uploadState(cfg, bucketID, stateDbContents)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotAuthenticated):
			log.Fatal(t(NotAuthenticated))
		default:
			log.Fatal(fmt.Errorf("uploading state: %w", err))
		}
	}

	err = deleteLocalStateDB(clusterManifest)
	if err != nil {
		log.Fatal(fmt.Errorf("deleting local state database: %w", err))
	}
}
