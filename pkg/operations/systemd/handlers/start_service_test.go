package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/systemd/handlers"
)

func TestNewStartService(t *testing.T) {
	tests := []struct {
		name      string
		mock      []exec.InterfaceCommandContextExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "systemctl",
						Args:        []string{"is-active", "dummy.service"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte{},
									Err:    nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "systemctl",
						Args:        []string{"start", "dummy.service"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: nil,
								},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: check cmd returns active",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "systemctl",
						Args:        []string{"is-active", "dummy.service"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte("active"),
									Err:    nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "systemctl",
						Args:        []string{"start", "dummy.service"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: nil,
								},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: enabled fail",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "systemctl",
						Args:        []string{"is-active", "dummy.service"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte{},
									Err:    nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "systemctl",
						Args:        []string{"start", "dummy.service"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: xerrors.New("dummy"),
								},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectations(tt.mock)

			startService := handlers.NewStartService(e)
			err := startService.Run(context.Background(), "dummy.service")
			tt.errAssert(t, err)
		})
	}
}

func TestNewStartService__DryRun(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "success",
			in:   "dummy.service",
			out:  "systemctl start dummy.service\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})
			startService := handlers.NewStartService(nil)
			ctx := context.Background()
			ctx = ctxlib.SetDryRun(ctx, true)
			err := startService.Run(ctx, tt.in)
			assert.NoError(t, err)
			assert.Equal(t, tt.out, buf.String())
		})
	}
}
