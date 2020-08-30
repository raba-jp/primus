package cmd

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/exec"
	"github.com/raba-jp/primus/executor"
	"github.com/raba-jp/primus/functions"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
)

func NewApplyCommand(in io.Reader, out io.Writer, errout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply changes",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			fs := afero.NewOsFs()
			exec := exec.New()
			client := http.DefaultClient

			exc := executor.NewExecutorWithArgs(in, out, errout, exec, fs, client)

			predeclared := starlark.StringDict{
				"execute":      starlark.NewBuiltin("execute", functions.Command(ctx, exc)),
				"symlink":      starlark.NewBuiltin("symlink", functions.Symlink(ctx, exc)),
				"http_request": starlark.NewBuiltin("http_request", functions.HttpRequest(ctx, exc)),
				"package":      starlark.NewBuiltin("package", functions.Package(ctx, exc)),
				"file_copy":    starlark.NewBuiltin("file_copy", functions.FileCopy(ctx, exc)),
				"file_move":    starlark.NewBuiltin("file_move", functions.FileMove(ctx, exc)),
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
