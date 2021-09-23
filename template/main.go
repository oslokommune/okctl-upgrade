package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type cmdFlags struct {
	debug bool
	force bool
}

func buildRootCommand() *cobra.Command {
	flags := cmdFlags{}

	var context Context

	cmd := &cobra.Command{
		PreRunE: func(_ *cobra.Command, args []string) error {
			context = newContext(flags)
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			if !flags.force {
				context.logger.Info("Simulating the upgrade, not doing any changes.")
			}

			err := upgrade(context)
			if err != nil {
				return fmt.Errorf("upgrade failed: %w", err)
			}

			return nil
		},
	}

	/*
	 * Flags supported. Expected behavior is as following:
	 *
	 * --debug:		Outputs extra output for debugging.
	 *
	 * --dry-run: 	If set to true, the upgrade will not make any changes, but only print what would be done, as if
	 * 				running a simulation.
	 *				If set to false, the upgrade will make actual changes.
	 *
	 * --confirm:	Skips all confirmation prompts, if any.
	 */
	cmd.PersistentFlags().BoolVarP(&flags.debug, "debug", "d", false, "Set this to enable debug output.")
	cmd.PersistentFlags().BoolVarP(&flags.debug, "dry-run", "n", true, "Don't actually do any changes, just show what would be done.")
	cmd.PersistentFlags().BoolVarP(&flags.force, "confirm", "c", false, "Set this to skip confirmation prompts.")

	return cmd
}
