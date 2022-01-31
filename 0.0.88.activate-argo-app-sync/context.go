package main

import (
	"github.com/oslokommune/okctl-upgrade/0.0.88.activate-argo-app-sync/pkg/logger"
)

type Context struct {
	logger logger.Logger
}

func newContext(flags cmdFlags) Context {
	var level logger.Level
	if flags.debug {
		level = logger.Debug
	} else {
		level = logger.Info
	}

	return Context{
		logger: logger.New(level),
	}
}
