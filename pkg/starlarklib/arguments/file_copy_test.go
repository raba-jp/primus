package arguments_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
)

func TestNewFileCopyArguments(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   *arguments.FileCopyArguments
		hasErr bool
	}{
		{
			name: "success: kwargs",
			data: `file_copy(src="/sym/src.txt", dest="/sym/dest.txt")`,
			want: &arguments.FileCopyArguments{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
				Perm: 0o777,
			},
			hasErr: false,
		},
		{
			name: "success: args",
			data: `file_copy("/sym/src.txt", "/sym/dest.txt")`,
			want: &arguments.FileCopyArguments{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
				Perm: 0o777,
			},
			hasErr: false,
		},
		{
			name: "success: with permission",
			data: `file_copy("/sym/src.txt", "/sym/dest.txt", 0o644)`,
			want: &arguments.FileCopyArguments{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
				Perm: 0o644,
			},
			hasErr: false,
		},
		{
			name: "error: too many arguments",
			data: `file_copy("/sym/src.txt", "/sym/dest.txt", 0o644, "many")`,
			want: &arguments.FileCopyArguments{
				Src:  "",
				Dest: "",
				Perm: 0o000,
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *arguments.FileCopyArguments

			predeclared := starlark.StringDict{
				"file_copy": starlark.NewBuiltin("file_copy", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var err error
					got, err = arguments.NewFileCopyArguments(b, args, kwargs)
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
