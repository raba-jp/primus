package os

import (
	std "path/filepath"

	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func GetCurrentPath() starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		v := starlark.GetCurrentFilePath(thread)
		return lib.String(v), nil
	}
}

func GetDir() starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		path := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "path", &path); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		v := std.Dir(path)
		return lib.String(v), nil
	}
}

func JoinPath() starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		list := &lib.List{}
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "paths", &list); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
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
		path := std.Join(values...)

		return lib.String(path), nil
	}
}
