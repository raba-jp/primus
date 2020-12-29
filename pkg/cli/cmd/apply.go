package cmd

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/args"
	"github.com/raba-jp/primus/pkg/starlark/exec"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
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
			log.Ctx(ctx).Info().Str("filepath", path).Msg("entrypoint")
			execFile := exec.Initialize()
			if err := execFile(ctx, false, path); err != nil {
				log.Ctx(ctx).Err(err).Msg("Failed to exec")
				return xerrors.Errorf(": %w", err)
			}

			return nil
		},
	}
}
