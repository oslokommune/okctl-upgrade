package memlimit

import (
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.93-increase-argocd-memlimit/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.93-increase-argocd-memlimit/pkg/lib/logger"
)

// Increaser increases argocd-repo-server's memory limit
type Increaser struct {
	flags cmdflags.Flags
	log   logger.Logger
}

// Upgrade upgrades the component
func (c Increaser) Upgrade() error {
	c.log.Info("Increasing memory limit on ArgoCD's deployment argocd-repo-server from 256Mi to 512Mi")

	if c.flags.DryRun {
		c.log.Info("Simulating some stuff")
	} else {
		c.log.Info("Doing some stuff")
	}

	return nil
}

func New(logger logger.Logger, flags cmdflags.Flags) Increaser {
	return Increaser{
		log:   logger,
		flags: flags,
	}
}
