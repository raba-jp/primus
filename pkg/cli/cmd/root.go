package cmd

import (
	"os"

	"github.com/raba-jp/primus/pkg/cli"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewPrimusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "primus",
		Short: "provisioning tool for local machine",
	}

	cmd.AddCommand(
		NewPlanCommand(),
		NewApplyCommand(),
		NewVersionCommand(),
	)
	AddLoggingFlag(cmd)

	return cmd
}

func Execute() {
	cmd := NewPrimusCommand()
	if err := cmd.Execute(); err != nil {
		zap.L().Error("Failed to execute", zap.Error(err))
		os.Exit(1)
	}
}

func AddLoggingFlag(cmd *cobra.Command) {
	var debugEnabled bool
	cmd.PersistentFlags().BoolVar(&debugEnabled, "debug", false, "Debug level output")

	cobra.OnInitialize(func() {
		if !debugEnabled {
			cli.EnableLogger()
		} else {
			cli.EnableDebugLogger()
		}
	})
}
