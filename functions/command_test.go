package functions_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_executor "github.com/raba-jp/primus/executor/mock"
	"github.com/raba-jp/primus/functions"
	"go.starlark.net/starlark"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		data     string
		wantCmd  string
		wantArgs []string
		hasErr   bool
	}{
		{
			data:     `command(name="echo", args=["hello", "world"])`,
			wantCmd:  "echo",
			wantArgs: []string{"hello", "world"},
			hasErr:   false,
		},
		{
			data:     `command(name="echo", args=[1])`,
			wantCmd:  "echo",
			wantArgs: []string{"1"},
			hasErr:   false,
		},
		{
			data:     `command(name="echo", args=[False, True])`,
			wantCmd:  "echo",
			wantArgs: []string{"false", "true"},
			hasErr:   false,
		},
		{
			data:     `command(name="echo", args=[1.111])`,
			wantCmd:  "",
			wantArgs: nil,
			hasErr:   true,
		},
		{
			data:     `command("echo", ["hello", "world"])`,
			wantCmd:  "echo",
			wantArgs: []string{"hello", "world"},
			hasErr:   false,
		},
		{
			data:     `command("echo")`,
			wantCmd:  "echo",
			wantArgs: []string{},
			hasErr:   false,
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
