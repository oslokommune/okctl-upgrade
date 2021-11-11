package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		printUsage(os.Stdout)

		os.Exit(1)
	}

	preflightCheck()

	log.Println("Parsing cluster manifest")

	clusterManifest, err := parseClusterManifest(os.Args[1])
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Fatal(t(ClusterManifestNotFound))
		}

		log.Fatal(fmt.Errorf("acquiring cluster name: %w", err))
	}

	printUpgradeInfoMessage(os.Stdout, clusterManifest)

	log.Println("Starting upgrade")

	doUpgrade(clusterManifest)

	log.Println("Upgrade complete")
}
