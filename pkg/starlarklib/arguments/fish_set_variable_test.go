package arguments_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
)

func TestNewFishSetVariableArguments(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   *arguments.FishSetVariableArguments
		hasErr bool
	}{
		{
			name: "success: kwargs",
			data: `fish_set_variable(name="GOPATH", value="$HOME/go", scope="universal", export=True)`,
			want: &arguments.FishSetVariableArguments{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  arguments.FishVariableUniversalScope,
				Export: true,
			},
			hasErr: false,
		},
		{
			name: "success: int kwargs",
			data: `fish_set_variable("GOPATH", "$HOME/go", "universal", True)`,
			want: &arguments.FishSetVariableArguments{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  arguments.FishVariableUniversalScope,
				Export: true,
			},
			hasErr: false,
		},
		{
			name: "error: unexpected scope",
			data: `fish_set_variable(name="GOPATH", value="$HOME/go", scope="dummy", export=True)`,
			want: &arguments.FishSetVariableArguments{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Export: true,
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *arguments.FishSetVariableArguments

			predeclared := starlark.StringDict{
				"fish_set_variable": starlark.NewBuiltin("fish_set_variable", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var err error
					got, err = arguments.NewFishSetVariableArguments(b, args, kwargs)
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
