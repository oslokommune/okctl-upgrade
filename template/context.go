package main

import (
	"github.com/oslokommune/okctl-upgrade/template/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/template/pkg/lib/logging"
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
