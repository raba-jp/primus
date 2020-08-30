package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewPrimusCommand() *cobra.Command {
	return NewPrimusCommandWithArgs(os.Stdin, os.Stdout, os.Stderr)
}

func NewPrimusCommandWithArgs(in io.Reader, out, errout io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "primus",
		Short: "provisioning tool for local machine",
	}

	cmd.AddCommand(
		NewApplyCommand(in, out, errout),
		NewVersionCommand(in, out, errout),
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
			enableLogger()
		} else {
			enableDebugLogger()
		}
	})
}
