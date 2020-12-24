package special

import (
	"io"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/xerrors"
)

func NewPrintContextFunction() starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		logger := ctxlib.Logger(ctx)

		key := ctxlib.PrevilegedAccessKey(ctx)

		logger.Debug("Print context", zap.String("key", key))
		return lib.None, nil
	}
}

func NewRequirePrevilegedAccessFunction(reader io.Reader) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		_, logger := ctxlib.LoggerWithNamespace(ctx, "require_previleged_access")

		ui.Printf("Root Password >>> ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		ctx = ctxlib.SetPrevilegedAccessKey(ctx, string(password))
		starlark.SetContext(ctx, thread)

		logger.Debug("Set previleged access key", zap.String("key", string(password)))

		return lib.None, nil
	}
}
