package cmd

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/args"
	"github.com/raba-jp/primus/pkg/starlark/exec"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type ApplyCommand *cobra.Command

func NewApplyCommand() ApplyCommand {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply changes",
		Args:  args.SingleFileArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			path, err := filepath.Abs(args[0])
			if err != nil {
				return xerrors.Errorf("Failed to get absolute path: %w", err)
			}
			zap.L().Info("entrypoint", zap.String("filepath", path))
			execFile := exec.Initialize()
			if err := execFile(ctx, false, path); err != nil {
				zap.L().Error("Failed to exec", zap.Error(err))
				return xerrors.Errorf(": %w", err)
			}

			return nil
		},
	}
}
