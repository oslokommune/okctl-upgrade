package grafana

import (
	"errors"
	"fmt"
	"os"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/oslokommune/okctl-upgrade/0.0.78/pkg/commonerrors"
	"github.com/oslokommune/okctl-upgrade/0.0.78/pkg/logger"
)

// Upgrader is a sample okctl component
type Upgrader struct {
	logger  logger.Logger
	dryRun  bool
	confirm bool
}

// Upgrade upgrades the component
//nolint:funlen
func (c Upgrader) Upgrade() error {
	c.logger.Info("Upgrading Grafana")

	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		return errors.New("missing required KUBECONFIG environment variable")
	}

	kubectlClient, err := acquireKubectlClient(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("acquiring kubectl client: %w", err)
	}

	err = preflight(c.logger, kubectlClient)
	if err != nil {
		if errors.Is(err, ErrNothingToDo) {
			c.logger.Info("Grafana is either missing or already upgraded. Ignoring upgrade.")

			return nil
		}

		return fmt.Errorf("running preflight checks: %w", err)
	}

	c.logger.Debug(fmt.Sprintf("Passed preflight test. Upgrading Grafana to %s", targetGrafanaVersion.String()))

	if !c.dryRun && !c.confirm {
		c.logger.Info(fmt.Sprintf("This will bump Grafana to %s", targetGrafanaVersion.String()))

		answer, err := c.askUser("Do you want to proceed?")
		if err != nil {
			return fmt.Errorf("prompting user: %w", err)
		}

		if !answer {
			return commonerrors.ErrUserAborted
		}
	}

	if c.dryRun {
		c.logger.Info("Simulating upgrade")

		err = patchGrafanaDeployment(c.logger, kubectlClient, c.dryRun)
		if err != nil {
			return fmt.Errorf("patching grafana deployment: %w", err)
		}
	} else {
		c.logger.Info("Patching Grafana")

		err = patchGrafanaDeployment(c.logger, kubectlClient, c.dryRun)
		if err != nil {
			return fmt.Errorf("patching grafana deployment: %w", err)
		}
	}

	err = postflight(c.logger, kubectlClient, c.dryRun)
	if err != nil {
		return fmt.Errorf("running postflight checks: %w", err)
	}

	c.logger.Info("Upgrading Grafana done!")

	return nil
}

func (c Upgrader) askUser(question string) (bool, error) {
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

type Opts struct {
	DryRun  bool
	Confirm bool
}

func New(logger logger.Logger, opts Opts) Upgrader {
	return Upgrader{
		logger:  logger,
		dryRun:  opts.DryRun,
		confirm: opts.Confirm,
	}
}
