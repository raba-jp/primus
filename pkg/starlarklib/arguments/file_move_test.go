package arguments_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
)

func TestNewFileMoveArguments(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   *arguments.FileMoveArguments
		hasErr bool
	}{
		{
			name: "success: kwargs",
			data: `file_move(src="/sym/src.txt", dest="/sym/dest.txt")`,
			want: &arguments.FileMoveArguments{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			hasErr: false,
		},
		{
			name: "success: args",
			data: `file_move("/sym/src.txt", "/sym/dest.txt")`,
			want: &arguments.FileMoveArguments{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			hasErr: false,
		},
		{
			name: "error: too many arguments",
			data: `file_move("/sym/src.txt", "/sym/dest.txt", "too many")`,
			want: &arguments.FileMoveArguments{
				Src:  "",
				Dest: "",
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *arguments.FileMoveArguments

			predeclared := starlark.StringDict{
				"file_move": starlark.NewBuiltin("file_move", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var err error
					got, err = arguments.NewFileMoveArguments(b, args, kwargs)
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
