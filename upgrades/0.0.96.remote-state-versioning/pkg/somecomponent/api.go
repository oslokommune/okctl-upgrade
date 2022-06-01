package somecomponent

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.96.remote-state-versioning/pkg/lib/cmdflags"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.96.remote-state-versioning/pkg/lib/commonerrors"
	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.96.remote-state-versioning/pkg/lib/logger"
)

// SomeComponent is a sample okctl component
type SomeComponent struct {
	flags cmdflags.Flags
	log   logger.Logger
}

// Upgrade upgrades the component
func (c SomeComponent) Upgrade() error {
	c.log.Info("Upgrading SomeComponent")

	c.log.Debug("SomeComponent is on version 0.5. Updating to 0.6")

	if !c.flags.DryRun && !c.flags.Confirm {
		c.log.Info("This will delete all logs.")

		answer, err := c.askUser("Do you want to continue?")
		if err != nil {
			return fmt.Errorf("prompting user: %w", err)
		}

		if !answer {
			return commonerrors.ErrUserAborted
		}
	}

	if c.flags.DryRun {
		c.log.Info("Simulating some stuff")
	} else {
		c.log.Info("Doing some stuff")
	}

	c.log.Info("Upgrading SomeComponent done!")

	return nil
}

func (c SomeComponent) askUser(question string) (bool, error) {
	answer := false
	prompt := &survey.Confirm{
		Message: question,
	}

	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return false, err
	}

	return answer, nil
}

func New(logger logger.Logger, flags cmdflags.Flags) SomeComponent {
	return SomeComponent{
		log:   logger,
		flags: flags,
	}
}
