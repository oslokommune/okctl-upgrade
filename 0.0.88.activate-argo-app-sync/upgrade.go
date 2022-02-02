package main

import (
	"fmt"
	"github.com/oslokommune/okctl-upgrade/0.0.88.activate-argo-app-sync/pkg/argocd"
	"github.com/oslokommune/okctl-upgrade/0.0.88.activate-argo-app-sync/pkg/kubectl/binary"
	"github.com/oslokommune/okctl-upgrade/0.0.88.activate-argo-app-sync/pkg/okctl"
)

func upgrade(upgradeContext Context, _ cmdFlags) error {
	o, err := okctl.InitializeOkctl()
	if err != nil {
		return fmt.Errorf("initializing okctl: %w", err)
	}

	kubectlClient := binary.New(binary.NewOpts{
		Logger:              upgradeContext.logger,
		Fs:                  upgradeContext.Fs,
		BinaryProvider:      o.BinariesProvider,
		CredentialsProvider: o.CredentialsProvider,
		Cluster:             *o.Declaration,
	})

	err = argocd.SetupApplicationsSync(argocd.SetupApplicationsSyncOpts{
		Fs:      upgradeContext.Fs,
		Cluster: *o.Declaration,
		Kubectl: kubectlClient,
	})
	if err != nil {
		return fmt.Errorf("activating application folder synchronization: %w", err)
	}

	err = argocd.MigrateExistingApplicationManifests(upgradeContext.Fs, *o.Declaration)
	if err != nil {
		return fmt.Errorf("migrating existing ArgoCD application manifests: %w", err)
	}

	return nil
}
