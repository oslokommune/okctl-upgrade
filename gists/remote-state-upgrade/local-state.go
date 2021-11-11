package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
)

const (
	defaultOutputPath      = "infrastructure"
	defaultStateDBFileName = "state.db"
)

func loadStateDB(clusterManifest Cluster) (io.Reader, error) {
	dbPath := path.Join(getOutputPath(clusterManifest), clusterManifest.Metadata.Name, defaultStateDBFileName)

	f, err := os.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	defer func() {
		_ = f.Close()
	}()

	buf := bytes.Buffer{}

	_, err = io.Copy(&buf, f)
	if err != nil {
		return nil, fmt.Errorf("buffering file: %w", err)
	}

	return &buf, nil
}

func deleteLocalStateDB(clusterManifest Cluster) error {
	dbPath := path.Join(getOutputPath(clusterManifest), clusterManifest.Metadata.Name, defaultStateDBFileName)

	return os.Remove(dbPath)
}

func getOutputPath(clusterManifest Cluster) string {
	if clusterManifest.Github.OutputPath == "" {
		return defaultOutputPath
	}

	return clusterManifest.Github.OutputPath
}
