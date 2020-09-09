package arguments

import (
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

var _ Arguments = (*SymlinkArguments)(nil)

type SymlinkArguments struct {
	Arguments
	Src  string
	Dest string
}

func NewSymlinkArguments(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (*SymlinkArguments, error) {
	a := SymlinkArguments{}
	err := a.Parse(b, args, kwargs)
	return &a, err
}

func (a *SymlinkArguments) Parse(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) error {
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "src", &a.Src, "dest", &a.Dest); err != nil {
		return xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	return nil
}
