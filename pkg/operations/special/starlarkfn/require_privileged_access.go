package starlarkfn

import (
	"io"

	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/starlark"
	"go.uber.org/zap"

	"github.com/raba-jp/primus/pkg/cli/ui"
	lib "go.starlark.net/starlark"
	"golang.org/x/crypto/ssh/terminal"
)

func RequirePrevilegedAccess(reader io.Reader) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		_, logger := ctxlib.LoggerWithNamespace(ctx, "require_previleged_access")

		ui.Printf("Root Password >>> ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		ctx = ctxlib.SetPrevilegedAccessKey(ctx, string(password))
		starlark.SetContext(thread, ctx)

		logger.Debug("Set previleged access key", zap.String("key", string(password)))

		return lib.None, nil
	}
}
