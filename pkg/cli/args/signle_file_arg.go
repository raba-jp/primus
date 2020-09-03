package args

import (
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func SingleFileArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return xerrors.New("requires a filepath")
	}
	return nil
}
