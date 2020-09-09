package arguments_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
)

func TestNewCommandArgs(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   *arguments.CommandArgs
		hasErr bool
	}{
		{
			name: "success: string array kwargs",
			data: `command(name="echo", args=["hello", "world"])`,
			want: &arguments.CommandArgs{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
				User: "",
				Cwd:  "",
			},
			hasErr: false,
		},
		{
			name: "success: int kwargs",
			data: `command(name="echo", args=[1])`,
			want: &arguments.CommandArgs{
				Cmd:  "echo",
				Args: []string{"1"},
				User: "",
				Cwd:  "",
			},
			hasErr: false,
		},
		{
			name: "error: bitint kwargs",
			data: `command(name="echo", args=[9007199254740991])`,
			want: &arguments.CommandArgs{
				Cmd: "echo",
			},
			hasErr: true,
		},
		{
			name: "success: bool kwargs",
			data: `command(name="echo", args=[False, True])`,
			want: &arguments.CommandArgs{Cmd: "echo",
				Args: []string{"false", "true"},
				User: "",
				Cwd:  "",
			},
			hasErr: false,
		},
		{
			name:   "success(unsupported): float kwargs",
			data:   `command(name="echo", args=[1.111])`,
			want:   nil,
			hasErr: true,
		},
		{
			name: "success: string args",
			data: `command("echo", ["hello", "world"])`,
			want: &arguments.CommandArgs{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
				User: "",
				Cwd:  "",
			},
			hasErr: false,
		},
		{
			name:   "success: no args",
			data:   `command("echo")`,
			want:   &arguments.CommandArgs{Cmd: "echo", Args: []string{}, User: "", Cwd: ""},
			hasErr: false,
		},
		{
			name: "success: set user and cwd",
			data: `command(name="echo", args=[], user="testuser", cwd="/home/testuser")`,
			want: &arguments.CommandArgs{
				Cmd:  "echo",
				Args: []string{},
				User: "testuser",
				Cwd:  "/home/testuser",
			},
			hasErr: false,
		},
		{
			name:   "error: invalid arg order",
			data:   `command("echo", "testuser", "/home/testuser", [])`,
			want:   &arguments.CommandArgs{Cmd: "echo"},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *arguments.CommandArgs

			predeclared := starlark.StringDict{
				"command": starlark.NewBuiltin("command", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var err error
					got, err = arguments.NewCommandArgs(b, args, kwargs)
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
