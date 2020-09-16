package cmd

import (
	"context"
	"os"

	"github.com/raba-jp/primus/pkg/cli/logging"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/spf13/cobra"
)

func NewCommand(planCmd PlanCommand, applyCmd ApplyCommand, versionCmd VersionCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "primus",
		Short: "provisioning tool for local machine",
	}

	cmd.AddCommand(planCmd, applyCmd, versionCmd)
	AddLoggingFlag(cmd)

	return cmd
}

func Execute() {
	cmd := Initialize()
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}

func AddLoggingFlag(cmd *cobra.Command) {
	var logLevel string

	cmd.PersistentFlags().StringVar(&logLevel, "logLevel", "", "Set log level. Allow info, debug, warn, and error")
	cobra.OnInitialize(func() {
		if err := logging.EnableLogger(logLevel); err != nil {
			ui.Errorf("%s", err)
			os.Exit(1)
		}
	})
}
