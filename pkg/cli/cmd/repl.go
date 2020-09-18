package cmd

import (
	"github.com/raba-jp/primus/pkg/internal/repl"
	"github.com/spf13/cobra"
)

type ReplCommand *cobra.Command

func NewReplCommand() ReplCommand {
	return &cobra.Command{
		Use:   "repl",
		Short: "Start REPL",
		RunE: func(cmd *cobra.Command, args []string) error {
			repl.NewPrompt()
			return nil
		},
	}
}
