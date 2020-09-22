package arguments_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
)

func TestNewCreateDirectoryArguments(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   *arguments.CreateDirectoryArguments
		hasErr bool
	}{
		{
			name: "success: kwargs",
			data: `create_directory(path="/sym/test", permission=0o777)`,
			want: &arguments.CreateDirectoryArguments{
				Path:       "/sym/test",
				Permission: 0o777,
			},
			hasErr: false,
		},
		{
			name: "success: args",
			data: `create_directory("/sym/test", 0o777)`,
			want: &arguments.CreateDirectoryArguments{
				Path:       "/sym/test",
				Permission: 0o777,
			},
			hasErr: false,
		},
		{
			name: "success: without permission",
			data: `create_directory(path="/sym/test")`,
			want: &arguments.CreateDirectoryArguments{
				Path:       "/sym/test",
				Permission: 0o644,
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `create_directory("/sym/test" 0o644, "too many")`,
			want:   nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *arguments.CreateDirectoryArguments

			predeclared := starlark.StringDict{
				"create_directory": starlark.NewBuiltin("create_directory", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var err error
					got, err = arguments.NewCreateDirectoryArguments(b, args, kwargs)
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
