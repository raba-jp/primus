package arguments

import (
	"os"

	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

var _ Arguments = (*CreateDirectoryArguments)(nil)

type CreateDirectoryArguments struct {
	Arguments
	Path       string
	Permission os.FileMode
}

func NewCreateDirectoryArguments(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (*CreateDirectoryArguments, error) {
	a := CreateDirectoryArguments{}
	err := a.Parse(b, args, kwargs)
	return &a, err
}

func (a *CreateDirectoryArguments) Parse(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) error {
	var perm = 0o644
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "path", &a.Path, "permission?", &perm); err != nil {
		return xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	a.Permission = os.FileMode(perm)
	return nil
}
