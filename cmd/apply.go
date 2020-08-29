package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/raba-jp/starlark_iac/functions"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"k8s.io/utils/exec"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fs := afero.NewOsFs()
		exec := exec.New()

		predeclared := starlark.StringDict{
			"execute": starlark.NewBuiltin("execute", functions.Execute(context.Background(), exec)),
			"symlink": starlark.NewBuiltin("symlink", functions.Symlink(context.Background(), fs)),
		}

		wd, _ := os.Getwd()
		dummy := filepath.Join(wd, "example.star")
		data, _ := afero.ReadFile(fs, dummy)
		thread := &starlark.Thread{
			Name: "apply",
		}
		_, err := starlark.ExecFile(thread, dummy, data, predeclared)
		if err != nil {
			zap.L().Error("Failed to exec", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
