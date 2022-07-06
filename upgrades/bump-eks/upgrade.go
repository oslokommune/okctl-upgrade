package main

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl-upgrade/upgrades/bump-eks/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/bump-eks/pkg/lib/logging"
	"github.com/spf13/afero"
)

func upgrade(ctx context.Context, log logging.Logger, fs *afero.Afero, flags cmdflags.Flags) error {
	clusterManifest, err := manifest.Cluster(fs)
	if err != nil {
		return fmt.Errorf("acquiring cluster manifest: %w", err)
	}

}
