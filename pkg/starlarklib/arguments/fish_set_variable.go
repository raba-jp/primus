package arguments

import (
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type FishVariableScope int

const (
	FishVariableUniversalScope FishVariableScope = iota + 1
	FishVariableGlobalScope
	FishVariableLocalScope
)

type FishSetVariableArguments struct {
	Arguments
	Name   string
	Value  string
	Scope  FishVariableScope
	Export bool
}

func NewFishSetVariableArguments(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (*FishSetVariableArguments, error) {
	a := FishSetVariableArguments{}
	err := a.Parse(b, args, kwargs)
	return &a, err
}

func (a *FishSetVariableArguments) Parse(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) error {
	var scope string
	if err := starlark.UnpackArgs(
		b.Name(),
		args,
		kwargs,
		"name", &a.Name,
		"value", &a.Value,
		"scope", &scope,
		"export", &a.Export,
	); err != nil {
		return xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	switch scope {
	case "universal":
		a.Scope = FishVariableUniversalScope
	case "global":
		a.Scope = FishVariableGlobalScope
	case "local":
		a.Scope = FishVariableLocalScope
	default:
		return xerrors.Errorf("Unexpected scope: %s", scope)
	}

	return nil
}
