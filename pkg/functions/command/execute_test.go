package command_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewExecuteFunction(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      command.ExecuteRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: string array kwargs",
			data: `test(cmd="echo", args=["hello", "world"])`,
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: int kwargs",
			data: `test(cmd="echo", args=[1])`,
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: bigint kwargs",
			data: `test(cmd="echo", args=[9007199254740991])`,
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "success: bool kwargs",
			data: `test(cmd="echo", args=[False, True])`,
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success(unsupported): float kwargs",
			data: `test(cmd="echo", args=[1.111])`,
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "success: no args",
			data: `test("echo")`,
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with user and cwd",
			data: `test("echo", [], user="testuser", cwd="/home/testuser")`,
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("echo", [], "testuser", "/home/testuser", "too many")`,
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: execute command failed",
			data: `test("echo")`,
			mock: func(ctx context.Context, p *command.Params) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := starlark.ExecForTest("test", tt.data, command.NewExecuteFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name       string
		params     *command.Params
		mock       []exec.InterfaceCommandContextExpectation
		mockStdout string
		mockErr    error
		errAssert  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			params: &command.Params{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
			},
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "echo",
						Args:        []string{"hello", "world"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
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
			name: "success: with user",
			params: &command.Params{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
				User: "root",
			},
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "echo",
						Args:        []string{"hello", "world"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplySetSysProcAttrExpectation(exec.CmdSetSysProcAttrExpectation{
								Args: exec.CmdSetSysProcAttrArgs{ProcAnything: true},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{Err: nil},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with cwd",
			params: &command.Params{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
				Cwd:  "/",
			},
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "echo",
						Args:        []string{"hello", "world"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplySetDirExpectation(exec.CmdSetDirExpectation{
								Args: exec.CmdSetDirArgs{Dir: "/"},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{Err: nil},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with user and cwd",
			params: &command.Params{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
				User: "root",
				Cwd:  "/",
			},
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "echo",
						Args:        []string{"hello", "world"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplySetDirExpectation(exec.CmdSetDirExpectation{
								Args: exec.CmdSetDirArgs{Dir: "/"},
							})
							cmd.ApplySetSysProcAttrExpectation(exec.CmdSetSysProcAttrExpectation{
								Args: exec.CmdSetSysProcAttrArgs{ProcAnything: true},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{Err: nil},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with stdin",
			params: &command.Params{
				Cmd:   "echo",
				Args:  []string{"hello", "world"},
				Stdin: new(bytes.Buffer),
			},
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "echo",
						Args:        []string{"hello", "world"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStdinExpectation(exec.CmdSetStdinExpectation{
								Args: exec.CmdSetStdinArgs{InAnything: true},
							})
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{Err: nil},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with stdout",
			params: &command.Params{
				Cmd:    "echo",
				Args:   []string{"hello", "world"},
				Stdout: new(bytes.Buffer),
			},
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "echo",
						Args:        []string{"hello", "world"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{Err: nil},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with stderr",
			params: &command.Params{
				Cmd:    "echo",
				Args:   []string{"hello", "world"},
				Stderr: new(bytes.Buffer),
			},
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "echo",
						Args:        []string{"hello", "world"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{Err: nil},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "failure",
			params: &command.Params{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
				User: "root",
				Cwd:  "/",
			},
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "echo",
						Args:        []string{"hello", "world"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplySetStdoutExpectation(exec.CmdSetStdoutExpectation{
								Args: exec.CmdSetStdoutArgs{OutAnything: true},
							})
							cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
								Args: exec.CmdSetStderrArgs{OutAnything: true},
							})
							cmd.ApplySetDirExpectation(exec.CmdSetDirExpectation{
								Args: exec.CmdSetDirArgs{Dir: "/"},
							})
							cmd.ApplySetSysProcAttrExpectation(exec.CmdSetSysProcAttrExpectation{
								Args: exec.CmdSetSysProcAttrArgs{ProcAnything: true},
							})
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{Err: xerrors.New("dummy")},
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
			exc := new(exec.MockInterface)
			exc.ApplyCommandContextExpectations(tt.mock)

			err := command.Execute(exc)(context.Background(), tt.params)
			tt.errAssert(t, err)
		})
	}
}
