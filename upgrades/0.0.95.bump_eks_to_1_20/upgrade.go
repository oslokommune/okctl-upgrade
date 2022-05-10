package main

import (
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.95.bump_eks_to_1_20/pkg/eks"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.95.bump_eks_to_1_20/pkg/lib/cmdflags"
)

func upgrade(context Context, flags cmdflags.Flags) error {
	c := eks.New(context.logger, flags)

	err := c.Upgrade()
	if err != nil {
		return err
	}

	return nil
}
