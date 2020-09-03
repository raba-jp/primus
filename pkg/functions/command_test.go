package functions_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/executor"
	mock_executor "github.com/raba-jp/primus/pkg/executor/mock"
	"github.com/raba-jp/primus/pkg/functions"
	"go.starlark.net/starlark"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		data   string
		want   *executor.CommandParams
		hasErr bool
	}{
		{
			data:   `command(name="echo", args=["hello", "world"])`,
			want:   &executor.CommandParams{CmdName: "echo", CmdArgs: []string{"hello", "world"}},
			hasErr: false,
		},
		{
			data:   `command(name="echo", args=[1])`,
			want:   &executor.CommandParams{CmdName: "echo", CmdArgs: []string{"1"}},
			hasErr: false,
		},
		{
			data:   `command(name="echo", args=[False, True])`,
			want:   &executor.CommandParams{CmdName: "echo", CmdArgs: []string{"false", "true"}},
			hasErr: false,
		},
		{
			data:   `command(name="echo", args=[1.111])`,
			want:   nil,
			hasErr: true,
		},
		{
			data:   `command("echo", ["hello", "world"])`,
			want:   &executor.CommandParams{CmdName: "echo", CmdArgs: []string{"hello", "world"}},
			hasErr: false,
		},
		{
			data:   `command("echo")`,
			want:   &executor.CommandParams{CmdName: "echo", CmdArgs: []string{}},
			hasErr: false,
		},
		{
			data:   `command("echo", [], user="testuser", cwd="/home/testuser")`,
			want:   &executor.CommandParams{CmdName: "echo", CmdArgs: []string{}, User: "testuser", Cwd: "/home/testuser"},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_executor.NewMockExecutor(ctrl)
			if !tt.hasErr {
				m.EXPECT().Command(gomock.Any(), gomock.Any()).Return(true, nil)
			}

			predeclared := starlark.StringDict{
				"command": starlark.NewBuiltin("command", functions.Command(context.Background(), m)),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if !tt.hasErr && err != nil {
				t.Error(err)
			}
		})
	}
}
