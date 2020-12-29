package special

import (
	"io"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/rs/zerolog/log"
	lib "go.starlark.net/starlark"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/xerrors"
)

func NewPrintContextFunction() starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		key := ctxlib.PrevilegedAccessKey(ctx)

		log.Ctx(ctx).Debug().Str("key", key).Msg("Print context")
		return lib.None, nil
	}
}

func NewRequirePrevilegedAccessFunction(reader io.Reader) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		ui.Printf("Root Password >>> ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		ctx = ctxlib.SetPrevilegedAccessKey(ctx, string(password))
		starlark.SetContext(ctx, thread)

		log.Ctx(ctx).Debug().Str("key", string(password)).Msg("set previleged access key")

		return lib.None, nil
	}
}
