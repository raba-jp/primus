//+build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"
)

func Initialize() *cobra.Command {
	wire.Build(
		NewCommand,
		NewPlanCommand,
		NewApplyCommand,
		NewVersionCommand,
	)
	return nil
}
