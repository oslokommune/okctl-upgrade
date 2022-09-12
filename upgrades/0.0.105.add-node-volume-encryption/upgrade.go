package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/eksctl"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/logging"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/manifest/apis/okctl.io/v1alpha1"
)

func upgrade(context Context, flags cmdflags.Flags, cluster v1alpha1.Cluster) error {
	log := context.logger

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

	log.Debug("Deleting original nodegroups")

	err = eksctl.DeleteNodeGroups(cluster.Metadata.Name, originalNodegroupNames, flags.DryRun)
	if err != nil {
		return fmt.Errorf("deleting nodegroups: %w", err)
	}

	log.Debug("Update complete")

	return nil
}

func configLogPrinter(log logging.Logger, cfg io.Reader) (io.Reader, error) {
	buf := bytes.Buffer{}

	tee := io.TeeReader(cfg, &buf)

	raw, err := io.ReadAll(tee)
	if err != nil {
		return nil, fmt.Errorf("buffering: %w", err)
	}

	log.Debug(string(raw))

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
