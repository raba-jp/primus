package arguments

import (
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type FishSetPathArguments struct {
	Arguments
	Values []string
}

func NewFishSetPathArguments(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (*FishSetPathArguments, error) {
	a := FishSetPathArguments{}
	err := a.Parse(b, args, kwargs)
	return &a, err
}

func (a *FishSetPathArguments) Parse(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) error {
	list := &starlark.List{}
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "values", &list); err != nil {
		return xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	values := make([]string, 0, list.Len())

	iter := list.Iterate()
	defer iter.Done()
	var item starlark.Value
	for iter.Next(&item) {
		str, ok := starlark.AsString(item)
		if !ok {
			continue
		}
		values = append(values, str)
	}
	a.Values = values

	return nil
}
