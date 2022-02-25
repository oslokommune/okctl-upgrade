package main

import (
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.93-increase-argocd-memlimit/pkg/lib/cmdflags"
)

func upgrade(context Context, flags cmdflags.Flags) error {
	c := memlimit.New(context.logger, flags)

	err := c.Upgrade()
	if err != nil {
		return err
	}

	return nil
}
