package cli

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/raba-jp/primus/executor/plan"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func NewPlanCommand(in io.Reader, out io.Writer, errout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Show provisioning plan",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return xerrors.New("requires a entrypoint filepath")
			}
			if len(args) > 1 {
				return xerrors.New("requires only one filepath")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			buf := new(bytes.Buffer)
			exc := plan.NewPlanExecutorWithArgs(buf)

			path, err := filepath.Abs(args[0])
			if err != nil {
				return xerrors.Errorf("Failed to get absolute path: %w", err)
			}
			zap.L().Info("entrypoint", zap.String("filepath", path))

			if err := ExecStarlarkFile(ctx, exc, path); err != nil {
				zap.L().Error("Failed to exec", zap.Error(err))
				return xerrors.Errorf(": %w", err)
			}

			if _, err := io.Copy(out, buf); err != nil {
				return xerrors.Errorf(": %w", err)
			}
			return nil
		},
	}
}
