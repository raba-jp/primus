package cli

import (
	"io"
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/executor/plan"
	"github.com/raba-jp/primus/functions"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
)

func NewPlanCommand(in io.Reader, out io.Writer, errout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Show provisioning plan",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			exc := plan.NewPlanExecutorWithArgs(out, errout)

			predeclared := starlark.StringDict{
				"execute":      starlark.NewBuiltin("execute", functions.Command(ctx, exc)),
				"symlink":      starlark.NewBuiltin("symlink", functions.Symlink(ctx, exc)),
				"http_request": starlark.NewBuiltin("http_request", functions.HTTPRequest(ctx, exc)),
				"package":      starlark.NewBuiltin("package", functions.Package(ctx, exc)),
				"file_copy":    starlark.NewBuiltin("file_copy", functions.FileCopy(ctx, exc)),
				"file_move":    starlark.NewBuiltin("file_move", functions.FileMove(ctx, exc)),
			}

			wd, _ := os.Getwd()
			dummy := filepath.Join(wd, "example.star")
			fs := afero.NewOsFs()
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
