package context

import (
	"context"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/logger"
)

type Context struct {
	Ctx    context.Context
	Logger logger.Logger
	Flags  cmdflags.Flags
}
