package arguments_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
)

func TestNewPackageArguments(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   *arguments.PackageArguments
		hasErr bool
	}{
		{
			name: "success: kwargs",
			data: `package(name="base-devel")`,
			want: &arguments.PackageArguments{
				Name: "base-devel",
			},
			hasErr: false,
		},
		{
			name: "success: args",
			data: `package("base-devel")`,
			want: &arguments.PackageArguments{
				Name: "base-devel",
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `package("base-devel", "too many")`,
			want:   &arguments.PackageArguments{},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *arguments.PackageArguments

			predeclared := starlark.StringDict{
				"package": starlark.NewBuiltin("package", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var err error
					got, err = arguments.NewPackageArguments(b, args, kwargs)
					return starlark.None, err
				}),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if !tt.hasErr && err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
