package cmd

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/functions"

	"github.com/raba-jp/primus/pkg/cli/args"
	"github.com/raba-jp/primus/pkg/starlark"
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
			log.Info().Str("filepath", path).Send()
			execFile := starlark.Initialize()
			if err := execFile(ctx, path, functions.New()); err != nil {
				log.Err(err).Msg("Failed to exec")
				return xerrors.Errorf(": %w", err)
			}

			return nil
		},
	}
}
