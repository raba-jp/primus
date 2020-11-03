package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/systemd/handlers"
)

func TestNewEnableService(t *testing.T) {
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
						Args:        []string{"is-enabled", "dummy.service"},
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
						Args:        []string{"enable", "dummy.service"},
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
			name: "success: check cmd returns enabled",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "systemctl",
						Args:        []string{"is-enabled", "dummy.service"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte("enabled\n"),
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
						Args:        []string{"enable", "dummy.service"},
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
						Args:        []string{"is-enabled", "dummy.service"},
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
						Args:        []string{"enable", "dummy.service"},
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

			enableService := handlers.NewEnableService(e)
			err := enableService.Run(context.Background(), "dummy.service")
			tt.errAssert(t, err)
		})
	}
}

func TestNewEnableService__DryRun(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "success",
			in:   "dummy.service",
			out:  "systemctl enable dummy.service\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			enableService := handlers.NewEnableService(nil)

			ctx := context.Background()
			ctx = ctxlib.SetDryRun(ctx, true)

			if err := enableService.Run(ctx, tt.in); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.out, buf.String()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
