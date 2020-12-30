package fish

import (
	"context"
	"fmt"
	"strings"

	lib "go.starlark.net/starlark"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
)

type SetPathParams struct {
	Values []string
}

type SetPathRunner func(ctx context.Context, p *SetPathParams) error

func NewSetPathFunction(runner SetPathRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params, err := parseSetPathArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseSetPathArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*SetPathParams, error) {
	a := &SetPathParams{}

	list := &lib.List{}
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "values", &list); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	values := make([]string, 0, list.Len())

	iter := list.Iterate()
	defer iter.Done()
	var item lib.Value
	for iter.Next(&item) {
		str, ok := lib.AsString(item)
		if !ok {
			continue
		}
		values = append(values, str)
	}
	a.Values = values

	return a, nil
}

func SetPath(execute backend.Execute) SetPathRunner {
	return func(ctx context.Context, p *SetPathParams) error {
		arg := fmt.Sprintf("'set --universal fish_user_paths %s'", strings.Join(p.Values, " "))

		if err := execute(ctx, &backend.ExecuteParams{
			Cmd:  "fish",
			Args: []string{"--command", arg},
		}); err != nil {
			return xerrors.Errorf("failed to set path: fish --command 'set --universal fish_user_path %s': %w", arg, err)
		}
		log.Ctx(ctx).Info().Strs("values", p.Values).Msg("set fish path")
		return nil
	}
}
