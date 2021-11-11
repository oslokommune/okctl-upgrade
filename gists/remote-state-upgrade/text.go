package main

import (
	"fmt"
	"io"
	"log"
	"text/template"
)

func printUsage(writer io.Writer) {
	_, _ = fmt.Fprint(writer, usage)
}

const usage = `
remote-state-upgrade moves the state database used to manage okctl state from your local machine to an s3 bucket.

Usage: ./remote-state-upgrade <path to cluster manifest>
Example: ./remote-state-upgrade cluster.yaml
`

func printUpgradeInfoMessage(writer io.Writer, clusterManifest Cluster) {
	upgradeInfoMessageTemplate, err := template.New("upgradeInfoMessage").Parse(upgradeInfoMessage)
	if err != nil {
		log.Fatal("parsing upgrade info message template: %w", err)
	}

	err = upgradeInfoMessageTemplate.Execute(writer, upgradeInfoMessageOpts{ClusterName: clusterManifest.Metadata.Name})
	if err != nil {
		log.Fatal(fmt.Errorf("printing info message: %w", err))
	}
}

type upgradeInfoMessageOpts struct {
	ClusterName string
}

const upgradeInfoMessage = `
This upgrade will move your state database to an S3 bucket. If anything goes wrong, you can revert the changes by:
1. Run 'git checkout -- .' to recover your local state database.
2. Go to S3 in the web console and empty the bucket called 'okctl-{{- .ClusterName -}}-meta'. The bucket needs to be empty so CloudFormation can delete it.
3. Go to CloudFormation in the web console and delete the stack called okctl-s3bucket-{{- .ClusterName -}}-okctl-{{- .ClusterName -}}-meta

`
