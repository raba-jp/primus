package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/operations/arch/handlers"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/arch/handlers/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewUninstall(t *testing.T) {
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
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: true},
			},
			mock: &exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "pacman",
					Args:        []string{"-R", "--noconfirm", "base-devel"},
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
			name: "success: not installed",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: false},
			},
			mock:      nil,
			errAssert: assert.NoError,
		},
		{
			name: "error: error occurred",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: true},
			},
			mock: &exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "pacman",
					Args:        []string{"-R", "--noconfirm", "base-devel"},
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

			uninstall := handlers.NewUninstall(handler, exc)
			err := uninstall.Run(context.Background(), false, &handlers.UninstallParams{Name: "base-devel"})
			tt.errAssert(t, err)
		})
	}
}

func TestNewUninstall__dryrun(t *testing.T) {
	tests := []struct {
		name   string
		params *handlers.UninstallParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.UninstallParams{
				Name: "pkg",
			},
			want: "pacman -R --noconfirm pkg\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			uninstall := handlers.NewUninstall(nil, nil)
			err := uninstall.Run(context.Background(), true, tt.params)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
