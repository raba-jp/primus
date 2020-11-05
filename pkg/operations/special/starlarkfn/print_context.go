package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
)

func PrintContext() starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		logger := ctxlib.Logger(ctx)

		key := ctxlib.PrevilegedAccessKey(ctx)
		dryrun := ctxlib.DryRun(ctx)

		logger.Debug("Print context", zap.String("key", key), zap.Bool("dryrun", dryrun))
		return lib.None, nil
	}
}
