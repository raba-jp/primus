package fish

import (
	"context"
	"fmt"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/rs/zerolog/log"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type SetVariableParams struct {
	Name   string
	Value  string
	Scope  VariableScope
	Export bool
}

type SetVariableRunner func(ctx context.Context, p *SetVariableParams) error

func NewSetVariableFunction(runner SetVariableRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params, err := parseSetVariableArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseSetVariableArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*SetVariableParams, error) {
	a := &SetVariableParams{}

	var scope string
	if err := lib.UnpackArgs(
		b.Name(),
		args,
		kwargs,
		"name", &a.Name,
		"value", &a.Value,
		"scope", &scope,
		"export", &a.Export,
	); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	switch scope {
	case "universal":
		a.Scope = UniversalScope
	case "global":
		a.Scope = GlobalScope
	case "local":
		a.Scope = LocalScope
	default:
		return nil, xerrors.Errorf("Unexpected scope: %s", scope)
	}

	return a, nil
}

func SetVariable(execute backend.Execute) SetVariableRunner {
	return func(ctx context.Context, p *SetVariableParams) error {
		var scope string
		switch p.Scope {
		case UniversalScope:
			scope = "--universal"
		case GlobalScope:
			scope = "--global"
		case LocalScope:
			scope = "--local"
		}

		export := ""
		if p.Export {
			export = " --export"
		}

		arg := fmt.Sprintf("'set %s%s %s %s'", scope, export, p.Name, p.Value)

		if err := execute(ctx, &backend.ExecuteParams{
			Cmd:  "fish",
			Args: []string{"--command", arg},
		}); err != nil {
			return xerrors.Errorf("failed to set variable: fish --command %s: %w", arg, err)
		}
		log.Info().
			Str("name", p.Name).
			Str("value", p.Value).
			Str("scope", scope).
			Bool("export", p.Export).
			Msg("set fish variable")

		return nil
	}
}
