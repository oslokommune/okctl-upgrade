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

	currentGrafanaVersion, err := getCurrentGrafanaVersion(kubectlClient)
	if err != nil {
		return fmt.Errorf("getting current Grafana version: %w", err)
	}

	err = validateVersion(expectedGrafanaVersionPreUpgrade, currentGrafanaVersion)
	if err != nil {
		return fmt.Errorf("unexpected Grafana version installed: %w", err)
	}

	c.logger.Debug(fmt.Sprintf(
		"Grafana is on version %s. Updating to %s",
		currentGrafanaVersion,
		upgradeTag,
	))

	if !c.dryRun && !c.confirm {
		c.logger.Info("This will delete all logs.")

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

		receipts, err := patchGrafanaDeployment(kubectlClient, c.dryRun)
		if err != nil {
			return fmt.Errorf("patching grafana deployment: %w", err)
		}

		for _, r := range receipts.receipts {
			c.logger.Info(r)
		}
	} else {
		c.logger.Info("Patching Grafana")

		_, err = patchGrafanaDeployment(kubectlClient, c.dryRun)
		if err != nil {
			return fmt.Errorf("patching grafana deployment: %w", err)
		}
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
