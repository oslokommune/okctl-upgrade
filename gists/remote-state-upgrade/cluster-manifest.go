package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"sigs.k8s.io/yaml"
)

type ClusterGithub struct {
	OutputPath string `json:"outputPath"`
}

type ClusterMetadata struct {
	Name string `json:"name"`
}

type Cluster struct {
	Metadata ClusterMetadata `json:"metadata"`
	Github   ClusterGithub   `json:"github"`
}

func parseClusterManifest(clusterManifestPath string) (Cluster, error) {
	f, err := os.Open(clusterManifestPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Cluster{}, ErrNotFound
		}

		log.Fatal(fmt.Errorf("opening cluster manifest: %w", err))
	}

	defer func() {
		_ = f.Close()
	}()

	raw, err := io.ReadAll(f)
	if err != nil {
		return Cluster{}, fmt.Errorf("reading data: %w", err)
	}

	payload := Cluster{}

	err = yaml.Unmarshal(raw, &payload)
	if err != nil {
		return Cluster{}, fmt.Errorf("unmarshalling data: %w", err)
	}

	return payload, nil
}
