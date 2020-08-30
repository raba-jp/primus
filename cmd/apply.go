package cmd

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/functions"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"k8s.io/utils/exec"
)

func NewApplyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply changes",
		Run: func(cmd *cobra.Command, args []string) {
			fs := afero.NewOsFs()
			exec := exec.New()

			predeclared := starlark.StringDict{
				"execute":      starlark.NewBuiltin("execute", functions.Execute(context.Background(), exec)),
				"symlink":      starlark.NewBuiltin("symlink", functions.Symlink(context.Background(), fs)),
				"http_request": starlark.NewBuiltin("http_request", functions.HttpRequest(context.Background(), http.DefaultClient, fs)),
				"package":      starlark.NewBuiltin("package", functions.Package(context.Background(), exec)),
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
}
