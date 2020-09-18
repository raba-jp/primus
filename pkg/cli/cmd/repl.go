package cmd

import "github.com/spf13/cobra"

type ReplCommand *cobra.Command

func NewReplCommand() ReplCommand {
	return &cobra.Command{
		Use:   "repl",
		Short: "Start REPL",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
