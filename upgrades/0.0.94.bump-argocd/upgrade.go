package main

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.bump-argocd/pkg/kubectl"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.bump-argocd/pkg/lib/cmdflags"
)

var (
	expectedVersion = *semver.MustParse("2.1.7")
	targetVersion   = *semver.MustParse("2.1.15")
)

func upgrade(context Context, flags cmdflags.Flags) error {
	argocdServerSelector := kubectl.Selector{
		Namespace:     "argocd",
		Kind:          "deployment",
		Name:          "argocd-server",
		ContainerName: "server",
	}

	currentVersion, err := kubectl.GetImageVersion(argocdServerSelector)
	if err != nil {
		return fmt.Errorf("acquiring argocd image version: %w", err)
	}

	if !expectedVersion.Equal(&currentVersion) {
		return fmt.Errorf("found version %s, ignoring", currentVersion.String())
	}

	fmt.Printf("Found version %s", currentVersion.String())

	err = kubectl.UpdateImageVersion(argocdServerSelector, targetVersion)
	if err != nil {
		return fmt.Errorf("updating version: %w", err)
	}

	return nil
}
