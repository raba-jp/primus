package cmd

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/spf13/cobra"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			ui.Printf("%s", Version)
		},
	}
}
