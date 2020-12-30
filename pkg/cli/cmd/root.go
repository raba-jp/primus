package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

func NewCommand(
	applyCmd ApplyCommand,
	versionCmd VersionCommand,
	replCmd ReplCommand,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "primus",
		Short: "provisioning tool for local machine",
	}

	cmd.AddCommand(applyCmd, versionCmd, replCmd)
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
	cobra.OnInitialize(func() {})
}
