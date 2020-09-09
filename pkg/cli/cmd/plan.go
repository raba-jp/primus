package cmd

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/args"
	"github.com/raba-jp/primus/pkg/executor/plan"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func NewPlanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Show provisioning plan",
		Args:  args.SingleFileArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			exc := plan.NewPlanExecutor()

			path, err := filepath.Abs(args[0])
			if err != nil {
				return xerrors.Errorf("Failed to get absolute path: %w", err)
			}
			zap.L().Info("entrypoint", zap.String("filepath", path))

			if err := functions.ExecStarlarkFile(ctx, exc, path); err != nil {
				zap.L().Error("Failed to exec", zap.Error(err))
				return xerrors.Errorf(": %w", err)
			}

			return nil
		},
	}
}
