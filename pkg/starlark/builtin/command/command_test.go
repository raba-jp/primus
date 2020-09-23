package command_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/handlers/mock"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/raba-jp/primus/pkg/starlark/builtin/command"
	"golang.org/x/xerrors"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockCommandHandler)
		hasErr bool
	}{
		{
			name: "success: string array kwargs",
			data: `test(name="echo", args=["hello", "world"])`,
			mock: func(m *mock_handlers.MockCommandHandler) {
				m.EXPECT().Command(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.CommandParams{
						CmdName: "echo",
						CmdArgs: []string{"hello", "world"},
						User:    "",
						Cwd:     "",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: int kwargs",
			data: `test(name="echo", args=[1])`,
			mock: func(m *mock_handlers.MockCommandHandler) {
				m.EXPECT().Command(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.CommandParams{
						CmdName: "echo",
						CmdArgs: []string{"1"},
						User:    "",
						Cwd:     "",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: bigint kwargs",
			data:   `test(name="echo", args=[9007199254740991])`,
			mock:   func(m *mock_handlers.MockCommandHandler) {},
			hasErr: true,
		},
		{
			name: "success: bool kwargs",
			data: `test(name="echo", args=[False, True])`,
			mock: func(m *mock_handlers.MockCommandHandler) {
				m.EXPECT().Command(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.CommandParams{
						CmdName: "echo",
						CmdArgs: []string{"false", "true"},
						User:    "",
						Cwd:     "",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "success(unsupported): float kwargs",
			data:   `test(name="echo", args=[1.111])`,
			mock:   func(m *mock_handlers.MockCommandHandler) {},
			hasErr: true,
		},
		{
			name: "success: no args",
			data: `test("echo")`,
			mock: func(m *mock_handlers.MockCommandHandler) {
				m.EXPECT().Command(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: with user and cwd",
			data: `test("echo", [], user="testuser", cwd="/home/testuser")`,
			mock: func(m *mock_handlers.MockCommandHandler) {
				m.EXPECT().Command(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("echo", [], "testuser", "/home/testuser", "too many")`,
			mock:   func(m *mock_handlers.MockCommandHandler) {},
			hasErr: true,
		},
		{
			name: "error: execute command failed",
			data: `test("echo")`,
			mock: func(m *mock_handlers.MockCommandHandler) {
				m.EXPECT().Command(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockCommandHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, command.Command(m))
			if !tt.hasErr && err != nil {
				t.Error(err)
			}
		})
	}
}
