package exec

import (
	"context"

	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type Fn = func(ctx context.Context, dryrun bool, path string) error

func NewExecFn(predeclared builtin.Predeclared, fs afero.Fs) Fn {
	return func(ctx context.Context, dryrun bool, path string) error {
		data, err := afero.ReadFile(fs, path)
		if err != nil {
			return xerrors.Errorf("Read file failed: %s: %w", path, err)
		}

		thread := starlark.NewThread(
			"main",
			starlark.WithLoad(starlark.Load(fs, predeclared)),
			starlark.WithContext(ctx),
			starlark.WithLogger(zap.L()),
			starlark.WithDryRunMode(dryrun),
		)
		if _, err := lib.ExecFile(thread, path, data, predeclared); err != nil {
			return xerrors.Errorf("Failed exec file: %w", err)
		}

		return nil
	}
}
