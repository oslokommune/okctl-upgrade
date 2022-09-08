package main

import (
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.105.add-node-volume-encryption/pkg/lib/logging"
)

type Context struct {
	logger logging.Logger
}

func newContext(flags cmdflags.Flags) Context {
	var level logging.Level
	if flags.Debug {
		level = logging.Debug
	} else {
		level = logging.Info
	}

	return Context{
		logger: logging.New(level),
	}
}
