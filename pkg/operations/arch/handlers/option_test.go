package handlers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/arch/handlers"
	"github.com/raba-jp/primus/pkg/operations/command/handlers/mocks"
)

func TestCmdArgs(t *testing.T) {
	tests := []struct {
		name    string
		mock    []mocks.ExecutableHandlerRunExpectation
		cmdType handlers.CmdType
		cmd     string
		args    []string
	}{
		{
			name: "install, yay, powerpill",
			mock: []mocks.ExecutableHandlerRunExpectation{
				{
					Args: mocks.ExecutableHandlerRunArgs{
						CtxAnything: true,
						Name:        "yay",
					},
					Returns: mocks.ExecutableHandlerRunReturns{Ok: true},
				},
				{
					Args: mocks.ExecutableHandlerRunArgs{
						CtxAnything: true,
						Name:        "powerpill",
					},
					Returns: mocks.ExecutableHandlerRunReturns{Ok: true},
				},
			},
			cmdType: handlers.InstallType,
			cmd:     "yay",
			args:    []string{"--pacman", "powerpill", "-S", "--noconfirm", "arg1", "arg2"},
		},
		{
			name: "install, yay",
			mock: []mocks.ExecutableHandlerRunExpectation{
				{
					Args: mocks.ExecutableHandlerRunArgs{
						CtxAnything: true,
						Name:        "yay",
					},
					Returns: mocks.ExecutableHandlerRunReturns{Ok: true},
				},
				{
					Args: mocks.ExecutableHandlerRunArgs{
						CtxAnything: true,
						Name:        "powerpill",
					},
					Returns: mocks.ExecutableHandlerRunReturns{Ok: false},
				},
			},
			cmdType: handlers.InstallType,
			cmd:     "yay",
			args:    []string{"-S", "--noconfirm", "arg1", "arg2"},
		},
		{
			name: "install, powerpill",
			mock: []mocks.ExecutableHandlerRunExpectation{
				{
					Args: mocks.ExecutableHandlerRunArgs{
						CtxAnything: true,
						Name:        "yay",
					},
					Returns: mocks.ExecutableHandlerRunReturns{Ok: false},
				},
				{
					Args: mocks.ExecutableHandlerRunArgs{
						CtxAnything: true,
						Name:        "powerpill",
					},
					Returns: mocks.ExecutableHandlerRunReturns{Ok: true},
				},
			},
			cmdType: handlers.InstallType,
			cmd:     "pacman",
			args:    []string{"-S", "--noconfirm", "arg1", "arg2"},
		},
		{
			name: "uninstall, yay",
			mock: []mocks.ExecutableHandlerRunExpectation{
				{
					Args: mocks.ExecutableHandlerRunArgs{
						CtxAnything: true,
						Name:        "yay",
					},
					Returns: mocks.ExecutableHandlerRunReturns{Ok: true},
				},
			},
			cmdType: handlers.UninstallType,
			cmd:     "yay",
			args:    []string{"-R", "--noconfirm", "arg1", "arg2"},
		},
		{
			name: "uninstall",
			mock: []mocks.ExecutableHandlerRunExpectation{
				{
					Args: mocks.ExecutableHandlerRunArgs{
						CtxAnything: true,
						Name:        "yay",
					},
					Returns: mocks.ExecutableHandlerRunReturns{Ok: false},
				},
			},
			cmdType: handlers.UninstallType,
			cmd:     "pacman",
			args:    []string{"-R", "--noconfirm", "arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executable := new(mocks.ExecutableHandler)
			executable.ApplyRunExpectations(tt.mock)

			cmd, args := handlers.CmdArgs(context.Background(), executable, tt.cmdType, []string{"arg1", "arg2"})
			assert.Equal(t, tt.cmd, cmd)
			assert.Equal(t, tt.args, args)
		})
	}
}
