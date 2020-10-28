package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/operations/arch/handlers"
	"github.com/raba-jp/primus/pkg/operations/arch/handlers/mocks"

	commandMock "github.com/raba-jp/primus/pkg/operations/command/handlers/mocks"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewInstall(t *testing.T) {
	checkInstallArgs := mocks.CheckInstallHandlerRunArgs{
		CtxAnything:  true,
		NameAnything: true,
	}

	executableArgs := commandMock.ExecutableHandlerRunArgs{
		CtxAnything:  true,
		NameAnything: true,
	}
	executables := []commandMock.ExecutableHandlerRunExpectation{
		{
			Args:    executableArgs,
			Returns: commandMock.ExecutableHandlerRunReturns{Ok: true},
		},
	}

	tests := []struct {
		name         string
		checkInstall mocks.CheckInstallHandlerRunExpectation
		executable   []commandMock.ExecutableHandlerRunExpectation
		mock         *exec.InterfaceCommandContextExpectation
		errAssert    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: false},
			},
			executable: executables,
			mock: &exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "yay",
					Args:        []string{"--pacman", "powerpill", "-S", "--noconfirm", "options", "base-devel"},
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyRunExpectation(exec.CmdRunExpectation{
							Returns: exec.CmdRunReturns{
								Err: nil,
							},
						})
						return cmd
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: already installed",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: true},
			},
			executable: executables,
			mock:       nil,
			errAssert:  assert.NoError,
		},
		{
			name: "error: install failed",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: false},
			},
			executable: executables,
			mock: &exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "yay",
					Args:        []string{"--pacman", "powerpill", "-S", "--noconfirm", "options", "base-devel"},
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyRunExpectation(exec.CmdRunExpectation{
							Returns: exec.CmdRunReturns{
								Err: xerrors.New("dummy"),
							},
						})
						return cmd
					},
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exc := new(exec.MockInterface)
			if tt.mock != nil {
				exc.ApplyCommandContextExpectation(*tt.mock)
			}

			checkInstall := new(mocks.CheckInstallHandler)
			checkInstall.ApplyRunExpectation(tt.checkInstall)

			executable := new(commandMock.ExecutableHandler)
			executable.ApplyRunExpectations(tt.executable)

			install := handlers.NewInstall(checkInstall, executable, exc)
			err := install.Run(context.Background(), false, &handlers.InstallParams{
				Name:   "base-devel",
				Option: "options",
			})
			tt.errAssert(t, err)
		})
	}
}

func TestNewInstall__dryrun(t *testing.T) {
	executableArgs := commandMock.ExecutableHandlerRunArgs{
		CtxAnything:  true,
		NameAnything: true,
	}
	executables := []commandMock.ExecutableHandlerRunExpectation{
		{
			Args:    executableArgs,
			Returns: commandMock.ExecutableHandlerRunReturns{Ok: true},
		},
		{
			Args:    executableArgs,
			Returns: commandMock.ExecutableHandlerRunReturns{Ok: true},
		},
	}

	tests := []struct {
		name   string
		params *handlers.InstallParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.InstallParams{
				Name:   "pkg",
				Option: "option",
			},
			want: "yay --pacman powerpill -S --noconfirm option pkg\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			executable := new(commandMock.ExecutableHandler)
			executable.ApplyRunExpectations(executables)

			install := handlers.NewInstall(nil, executable, nil)
			err := install.Run(context.Background(), true, tt.params)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
