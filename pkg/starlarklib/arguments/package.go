package arguments

import (
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

var _ Arguments = (*PackageArguments)(nil)

type PackageArguments struct {
	Arguments
	Name string
}

func NewPackageArguments(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (*PackageArguments, error) {
	a := PackageArguments{}
	err := a.Parse(b, args, kwargs)
	return &a, err
}

func (a *PackageArguments) Parse(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) error {
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "name", &a.Name); err != nil {
		return xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	return nil
}
