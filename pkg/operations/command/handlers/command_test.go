package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"golang.org/x/xerrors"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/command/handlers"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name       string
		params     *handlers.CommandParams
		mock       []exec.InterfaceCommandContextExpectation
		mockStdout string
		mockErr    error
		errAssert  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
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
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				User:    "root",
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
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				Cwd:     "/",
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
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				User:    "root",
				Cwd:     "/",
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
			name: "error: ",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				User:    "root",
				Cwd:     "/",
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
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectations(tt.mock)

			handler := handlers.NewCommand(e)

			err := handler.Command(context.Background(), false, tt.params)
			tt.errAssert(t, err)
		})
	}
}

func TestNewCommand__DryRun(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		want    string
	}{
		{
			name:    "no args",
			command: "ls",
			args:    []string{},
			want:    "ls \n",
		},
		{
			name:    "add option",
			command: "ls",
			args:    []string{"-al"},
			want:    "ls -al \n",
		},
		{
			name:    "with double quote",
			command: "ls",
			args:    []string{"-al", "\"go.mod\""},
			want:    "ls -al \"go.mod\" \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.NewCommand(nil)
			err := handler.Command(context.Background(), true, &handlers.CommandParams{
				CmdName: tt.command,
				CmdArgs: tt.args,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
