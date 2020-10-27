package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/operations/arch/handlers"
	"github.com/raba-jp/primus/pkg/operations/arch/handlers/mocks"

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

	tests := []struct {
		name         string
		checkInstall mocks.CheckInstallHandlerRunExpectation
		mock         *exec.InterfaceCommandContextExpectation
		errAssert    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: false},
			},
			mock: &exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "pacman",
					Args:        []string{"-S", "--noconfirm", "options", "base-devel"},
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
			mock:      nil,
			errAssert: assert.NoError,
		},
		{
			name: "error: install failed",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: false},
			},
			mock: &exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "pacman",
					Args:        []string{"-S", "--noconfirm", "options", "base-devel"},
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

			handler := new(mocks.CheckInstallHandler)
			handler.ApplyRunExpectation(tt.checkInstall)

			install := handlers.NewInstall(handler, exc)
			err := install.Run(context.Background(), false, &handlers.InstallParams{
				Name:   "base-devel",
				Option: "options",
			})
			tt.errAssert(t, err)
		})
	}
}

func TestNewInstall__dryrun(t *testing.T) {
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
			want: "pacman -S --noconfirm option pkg\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.NewInstall(nil, nil)
			err := handler.Run(context.Background(), true, tt.params)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
