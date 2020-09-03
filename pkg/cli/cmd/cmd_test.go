package cmd_test

import (
	"os"

	"github.com/spf13/cobra"
)

func executeCommand(child *cobra.Command, args ...string) error {
	p := cobra.Command{}
	p.AddCommand(child)
	p.SetOut(os.Stdout)
	p.SetErr(os.Stderr)

	a := append([]string{}, child.Use)
	a = append(a, args...)
	p.SetArgs(a)

	return p.Execute()
}
