package main

import (
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/somecomponent"
)

func upgrade(context Context, flags cmdflags.Flags) error {
	c := somecomponent.New(context.logger, flags)

	err := c.Upgrade()
	if err != nil {
		return err
	}

	return nil
}
