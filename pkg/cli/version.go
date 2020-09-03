package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func NewVersionCommand(in io.Reader, out, errout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(out, Version)
		},
	}
}
