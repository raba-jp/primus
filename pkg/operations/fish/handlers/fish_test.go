package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/fish/handlers"
	"golang.org/x/xerrors"
)

func TestNewSetVariable(t *testing.T) {
	tests := []struct {
		name      string
		params    *handlers.SetVariableParams
		mock      exec.InterfaceCommandContextExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: scope universal",
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.UniversalScope,
				Export: true,
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "fish",
					Args:        []string{"--command", "'set --universal --export GOPATH $HOME/go'"},
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
			errAssert: assert.NoError,
		},
		{
			name: "success: scope global",
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.GlobalScope,
				Export: true,
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "fish",
					Args:        []string{"--command", "'set --global --export GOPATH $HOME/go'"},
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
			errAssert: assert.NoError,
		},
		{
			name: "success: scope local",
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.LocalScope,
				Export: true,
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "fish",
					Args:        []string{"--command", "'set --local --export GOPATH $HOME/go'"},
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
			errAssert: assert.NoError,
		},
		{
			name: "success: no export",
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.LocalScope,
				Export: false,
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "fish",
					Args:        []string{"--command", "'set --local GOPATH $HOME/go'"},
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
			errAssert: assert.NoError,
		},
		{
			name: "error",
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.UniversalScope,
				Export: true,
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "fish",
					Args:        []string{"--command", "'set --universal --export GOPATH $HOME/go'"},
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
			exc.ApplyCommandContextExpectation(tt.mock)

			setVariable := handlers.NewSetVariable(exc)
			err := setVariable.Run(context.Background(), tt.params)
			tt.errAssert(t, err)
		})
	}
}

func TestNewSetVariable__DryRun(t *testing.T) {
	tests := []struct {
		name   string
		params *handlers.SetVariableParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.UniversalScope,
				Export: true,
			},
			want: "fish --command 'set --universal --export GOPATH $HOME/go'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			setVariable := handlers.NewSetVariable(nil)
			ctx := ctxlib.SetDryRun(context.Background(), true)
			err := setVariable.Run(ctx, tt.params)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestBaseBackend_FishSetPath(t *testing.T) {
	tests := []struct {
		name      string
		params    *handlers.SetPathParams
		mock      exec.InterfaceCommandContextExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			params: &handlers.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "fish",
					Args: []string{
						"--command",
						"'set --universal fish_user_paths $GOPATH/bin $HOME/.bin'",
					},
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
			errAssert: assert.NoError,
		},
		{
			name: "error",
			params: &handlers.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "fish",
					Args: []string{
						"--command",
						"'set --universal fish_user_paths $GOPATH/bin $HOME/.bin'",
					},
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
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectation(tt.mock)

			setPath := handlers.NewSetPath(e)
			err := setPath.Run(context.Background(), tt.params)
			tt.errAssert(t, err)
		})
	}
}

func TestNewSetPath__DryRun(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		params *handlers.SetPathParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			want: "fish --command 'set --universal fish_user_paths $GOPATH/bin $HOME/.bin'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			setPath := handlers.NewSetPath(nil)
			ctx := ctxlib.SetDryRun(context.Background(), true)
			err := setPath.Run(ctx, tt.params)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
