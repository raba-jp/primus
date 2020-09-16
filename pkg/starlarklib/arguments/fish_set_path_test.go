package arguments_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
)

func TestNewFishSetPathArguments(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   *arguments.FishSetPathArguments
		hasErr bool
	}{
		{
			name: "success: kwargs",
			data: `fish_set_path(values=["$GOPATH/bin", "$HOME/.bin"])`,
			want: &arguments.FishSetPathArguments{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			hasErr: false,
		},
		{
			name: "success: args",
			data: `fish_set_path(["$GOPATH/bin", "$HOME/.bin"])`,
			want: &arguments.FishSetPathArguments{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			hasErr: false,
		},
		{
			name: "success: include int and bool",
			data: `fish_set_path(["$GOPATH/bin", 1, True, "$HOME/.bin"])`,
			want: &arguments.FishSetPathArguments{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			hasErr: false,
		},
		{
			name: "error: too many argument",
			data: `fish_set_path(["$GOPATH/bin", "$HOME/.bin"], "too many")`,
			want: &arguments.FishSetPathArguments{
				Values: nil,
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *arguments.FishSetPathArguments

			predeclared := starlark.StringDict{
				"fish_set_path": starlark.NewBuiltin("fish_set_path", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var err error
					got, err = arguments.NewFishSetPathArguments(b, args, kwargs)
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
