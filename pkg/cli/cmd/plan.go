package cmd

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/args"
	"github.com/raba-jp/primus/pkg/starlark/exec"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

type PlanCommand *cobra.Command

func NewPlanCommand() PlanCommand {
	return &cobra.Command{
		Use:   "plan",
		Short: "Show provisioning plan",
		Args:  args.SingleFileArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			path, err := filepath.Abs(args[0])
			if err != nil {
				return xerrors.Errorf("Failed to get absolute path: %w", err)
			}
			log.Info().Str("filepath", path).Msg("entrypoint")

			execFile := exec.Initialize()
			if err := execFile(ctx, true, path); err != nil {
				log.Err(err).Msg("exec failed")
				return xerrors.Errorf(": %w", err)
			}

			return nil
		},
	}
}
