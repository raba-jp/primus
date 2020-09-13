package cmd

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/args"
	"github.com/raba-jp/primus/pkg/internal/backend"
	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func NewApplyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply changes",
		Args:  args.SingleFileArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			be := backend.New(exec.New(), afero.NewOsFs())

			path, err := filepath.Abs(args[0])
			if err != nil {
				return xerrors.Errorf("Failed to get absolute path: %w", err)
			}
			zap.L().Info("entrypoint", zap.String("filepath", path))

			if err := functions.ExecStarlarkFile(ctx, be, path); err != nil {
				zap.L().Error("Failed to exec", zap.Error(err))
				return xerrors.Errorf(": %w", err)
			}

			return nil
		},
	}
}
