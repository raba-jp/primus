package cmd

import (
	"github.com/raba-jp/primus/pkg/internal/promptlib"
	"github.com/spf13/cobra"
)

type ReplCommand *cobra.Command

func NewReplCommand() ReplCommand {
	return &cobra.Command{
		Use:   "repl",
		Short: "Start REPL",
		RunE: func(cmd *cobra.Command, args []string) error {
			promptlib.NewPrompt()
			return nil
		},
	}
}
