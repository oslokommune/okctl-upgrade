package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print(usage)

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

	upgradeInfoMessageTemplate, err := template.New("upgradeInfoMessage").Parse(upgradeInfoMessage)
	if err != nil {
		log.Fatal("parsing upgrade info message template: %w", err)
	}

	err = upgradeInfoMessageTemplate.Execute(os.Stdout, upgradeInfoMessageOpts{ClusterName: clusterManifest.Metadata.Name})
	if err != nil {
		log.Fatal(fmt.Errorf("printing info message: %w", err))
	}

	log.Println("Starting upgrade")

	doUpgrade(clusterManifest)

	log.Println("Upgrade complete")
}

const usage = `
remote-state-upgrade moves the state database used to manage okctl state from your local machine to an s3 bucket.

Usage: ./remote-state-upgrade <path to cluster manifest>
Example: ./remote-state-upgrade cluster.yaml
`

type upgradeInfoMessageOpts struct {
	ClusterName string
}

const upgradeInfoMessage = `
This upgrade will move your state database to an S3 bucket. If anything goes wrong, you can revert the changes by:
1. Run 'git checkout -- .' to recover your local state database.
2. Go to S3 in the web console and empty the bucket called 'okctl-{{- .ClusterName -}}-meta'. The bucket needs to be empty so CloudFormation can delete it.
3. Go to CloudFormation in the web console and delete the stack called okctl-s3bucket-{{- .ClusterName -}}-okctl-{{- .ClusterName -}}-meta

`