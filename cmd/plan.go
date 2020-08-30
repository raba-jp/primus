package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewPlanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Show provisioning plan",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("plan called")
		},
	}
}
