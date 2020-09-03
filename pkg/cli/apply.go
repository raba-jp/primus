package cli

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/executor/apply"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type NopWriter struct {
	io.Writer
}

func (*NopWriter) Write(buf []byte) (int, error) {
	return len(buf), nil
}

func NewApplyCommand(in io.Reader, out io.Writer, errout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply changes",
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

			fs := afero.NewOsFs()
			exec := exec.New()
			client := http.DefaultClient

			nop := new(NopWriter)

			exc := apply.NewApplyExecutorWithArgs(in, nop, errout, exec, fs, client)

			path, err := filepath.Abs(args[0])
			if err != nil {
				return xerrors.Errorf("Failed to get absolute path: %w", err)
			}
			zap.L().Info("entrypoint", zap.String("filepath", path))

			if err := ExecStarlarkFile(ctx, exc, path); err != nil {
				zap.L().Error("Failed to exec", zap.Error(err))
				return xerrors.Errorf(": %w", err)
			}

			return nil
		},
	}
}
