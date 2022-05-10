package eks

import (
	"errors"
	"fmt"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.95.bump_eks_to_1_20/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.95.bump_eks_to_1_20/pkg/lib/eks_upgrade"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.95.bump_eks_to_1_20/pkg/lib/logger"
)

// EKS is a sample okctl component
type EKS struct {
	flags  cmdflags.Flags
	log    logger.Logger
	eksctl eks_upgrade.Eks
}

// Upgrade upgrades the component
func (e EKS) Upgrade() error {
	e.log.Info("Upgrading EKS")

	err := e.preflight()
	if err != nil {
		if errors.Is(err, errNothingToDo) {
			return nil
		}

		return fmt.Errorf("running preflight checks: %w", err)
	}

	if e.flags.DryRun {
		e.log.Info("Simulating some stuff")
	} else {
		e.log.Info("Doing some stuff")
	}

	e.log.Info("Upgrading EKS done!")

	return nil
}

func (e EKS) preflight() error {
	version, err = e.eksctl.Version()
	if err != nil {
		return fmt.Errorf("getting version of EKS cluster: %w", err)
	}

	if version != eksSourceVersion {
		e.log.Info("EKS cluster is not on version 1.19, ignoring upgrade.")

		return errNothingToDo
	}

	e.log.Info("EKS cluster versino is %s, upgrading to %s", eksSourceVersion, eksTargetVersion)

	return nil
}

func New(logger logger.Logger, flags cmdflags.Flags) EKS {
	return EKS{
		log:    logger,
		flags:  flags,
		eksctl: eks_upgrade.New(),
	}
}
