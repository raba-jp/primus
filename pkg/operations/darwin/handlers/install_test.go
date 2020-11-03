package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/darwin/handlers"
	"github.com/raba-jp/primus/pkg/operations/darwin/handlers/mocks"
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
		mock         exec.InterfaceCommandContextExpectation
		params       *handlers.InstallParams
		errAssert    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: false},
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "brew",
					Args:        []string{"install", "options", "pkg"},
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
			params: &handlers.InstallParams{
				Name:   "pkg",
				Option: "options",
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: already installed",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: true},
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything:  true,
					CmdAnything:  true,
					ArgsAnything: true,
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
			params: &handlers.InstallParams{
				Name:   "pkg",
				Option: "options",
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: install package failed",
			checkInstall: mocks.CheckInstallHandlerRunExpectation{
				Args:    checkInstallArgs,
				Returns: mocks.CheckInstallHandlerRunReturns{Ok: false},
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "brew",
					Args:        []string{"install", "options", "pkg"},
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
			params: &handlers.InstallParams{
				Name:   "pkg",
				Option: "options",
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exc := new(exec.MockInterface)
			exc.ApplyCommandContextExpectation(tt.mock)

			handler := new(mocks.CheckInstallHandler)
			handler.ApplyRunExpectation(tt.checkInstall)

			install := handlers.NewInstall(handler, exc)
			err := install.Run(context.Background(), tt.params)
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
			want: "brew install option pkg\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			install := handlers.NewInstall(nil, nil)
			ctx := ctxlib.SetDryRun(context.Background(), true)
			err := install.Run(ctx, tt.params)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
