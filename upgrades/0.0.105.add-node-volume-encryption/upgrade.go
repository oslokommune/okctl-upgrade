package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/eksctl"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/logging"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/manifest/apis/okctl.io/v1alpha1"
)

func upgrade(context Context, flags cmdflags.Flags, cluster v1alpha1.Cluster) error {
	log := context.logger

	if !c.flags.DryRun && !c.flags.Confirm {
		continueUpgrade, err := prompt(continuationMessage)
		if err != nil {
			return fmt.Errorf("prompting for continuation: %w", err)
		}

		if !continueUpgrade {
			return errors.New("aborted by user")
		}
	}

	log.Debug("Retrieving nodegroup names")

	originalNodegroupNames, err := eksctl.GetNodeGroupNames(cluster.Metadata.Name)
	if err != nil {
		return fmt.Errorf("acquiring nodegroup names: %w", err)
	}

	log.Debugf("Found nodegroups %v\n", originalNodegroupNames)

	clusterVersion, err := eksctl.GetClusterVersion(cluster.Metadata.Name)
	if err != nil {
		return fmt.Errorf("acquiring cluster version: %w", err)
	}

	log.Debugf("Found cluster version %s\n", clusterVersion)

	log.Debug("Generating eksctl cluster configuration")

	cfg, err := eksctl.GenerateClusterConfig(
		cluster,
		generateNodegroupNames(cluster.Metadata.Region, clusterVersion, generateRandomString),
	)
	if err != nil {
		return fmt.Errorf("generating cluster config: %w", err)
	}

	cfg, err = configLogPrinter(log, cfg)
	if err != nil {
		return fmt.Errorf("debug printing cluster config: %w", err)
	}

	log.Debug("Updating nodegroups")

	err = eksctl.UpdateNodeGroups(cfg, flags.DryRun)
	if err != nil {
		return fmt.Errorf("updating nodegroups: %w", err)
	}

	err = deleteNodeGroups(cluster.Metadata.Name, originalNodegroupNames, flags.DryRun)
	if err != nil {
		return fmt.Errorf("deleting nodegroups: %w", err)
	}

	log.Debug("Deleting original nodegroups")

	log.Debug("Update complete")

	return nil
}

func deleteNodeGroups(clusterName string, nodegroupNames []string, dryRun bool) error {
	if len(nodegroupNames) == 0 {
		return nil
	}

	err := eksctl.DeleteNodeGroups(clusterName, nodegroupNames, dryRun)
	if err != nil {
		return err
	}

	return nil
}

func configLogPrinter(log logging.Logger, cfg io.Reader) (io.Reader, error) {
	buf := bytes.Buffer{}

	tee := io.TeeReader(cfg, &buf)

	raw, err := io.ReadAll(tee)
	if err != nil {
		return nil, fmt.Errorf("buffering: %w", err)
	}

	log.Debugf("\n%s", string(raw))

	return &buf, nil
}

func generateRandomString() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func generateNodegroupNames(region string, clusterVersion string, randomizerFn func() string) []string {
	regionParts := strings.Split(region, "-")
	availabilityZones := []string{"a", "b", "c"}
	availabilityZonePrefix := regionParts[len(regionParts)-1]
	postfix := strings.ToUpper(randomizerFn()[0:10])

	nodegroupNames := make([]string, len(availabilityZones))

	for index, az := range availabilityZones {
		nodegroupNames[index] = fmt.Sprintf(
			"ng-generic-%s-%s%s-%s",
			strings.ReplaceAll(clusterVersion, ".", "-"),
			availabilityZonePrefix,
			az,
			postfix,
		)
	}

	return nodegroupNames
}

func prompt(msg string) (bool, error) {
	answer := false
	prompt := &survey.Confirm{Message: msg}

	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return false, fmt.Errorf("prompting user: %w", err)
	}

	return answer, nil
}

const continuationMessage = `This upgrade will move applications to new nodes in the cluster, which can result in applications
experiencing downtime.

To ensure uptime, follow this guide: https://github.com/oslokommune/okctl-upgrade/tree/main/gists/bump-eks-to-1-20#alternative-1-no-downtime

Are you sure you want to continue with the upgrade?
`
