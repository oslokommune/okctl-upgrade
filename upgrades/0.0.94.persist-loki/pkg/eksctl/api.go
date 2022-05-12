package eksctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.94.persist-loki/pkg/lib/context"
)

func CreateServiceUser(ctx context.Context, clusterName string, name string, policies []string) error {
	eksctlPath, err := acquireEksctlPath(ctx.Fs, os.UserHomeDir)
	if err != nil {
		return fmt.Errorf("acquiring eksctl path: %w", err)
	}

	args := []string{
		"create", "iamserviceaccount",
		"--name", name,
		"--namespace", defaultMonitoringNamespace,
		"--cluster", clusterName,
		"--role-name", fmt.Sprintf("okctl-%s-loki", clusterName),
		"--attach-policy-arn", strings.Join(policies, ","),
		"--approve",
		"--override-existing-serviceaccounts",
	}

	if ctx.Flags.DryRun {
		ctx.Logger.Debugf("Running eksctl with args: %v", args)

		return nil
	}

	err = runEksctlCommand(eksctlPath, args...)
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}
