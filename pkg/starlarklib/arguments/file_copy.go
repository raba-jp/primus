package arguments

import (
	"os"

	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

var _ Arguments = (*FileCopyArguments)(nil)

type FileCopyArguments struct {
	Arguments
	Src  string
	Dest string
	Perm os.FileMode
}

func NewFileCopyArguments(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (*FileCopyArguments, error) {
	a := FileCopyArguments{}
	err := a.Parse(b, args, kwargs)
	return &a, err
}

func (a *FileCopyArguments) Parse(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) error {
	var perm = 0o777
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "src", &a.Src, "dest", &a.Dest, "permission?", &perm); err != nil {
		return xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	a.Perm = os.FileMode(perm)
	return nil
}
