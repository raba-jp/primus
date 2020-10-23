package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/operations/command/handlers"
	"github.com/raba-jp/primus/pkg/operations/command/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      handlers.CommandHandlerCommandExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: string array kwargs",
			data: `test(name="echo", args=["hello", "world"])`,
			mock: handlers.CommandHandlerCommandExpectation{
				Args: handlers.CommandHandlerCommandArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CommandParams{
						CmdName: "echo",
						CmdArgs: []string{"hello", "world"},
						User:    "",
						Cwd:     "",
					},
				},
				Returns: handlers.CommandHandlerCommandReturns{Err: nil},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: int kwargs",
			data: `test(name="echo", args=[1])`,
			mock: handlers.CommandHandlerCommandExpectation{
				Args: handlers.CommandHandlerCommandArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CommandParams{
						CmdName: "echo",
						CmdArgs: []string{"1"},
						User:    "",
						Cwd:     "",
					},
				},
				Returns: handlers.CommandHandlerCommandReturns{Err: nil},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: bigint kwargs",
			data:      `test(name="echo", args=[9007199254740991])`,
			mock:      handlers.CommandHandlerCommandExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "success: bool kwargs",
			data: `test(name="echo", args=[False, True])`,
			mock: handlers.CommandHandlerCommandExpectation{
				Args: handlers.CommandHandlerCommandArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CommandParams{
						CmdName: "echo",
						CmdArgs: []string{"false", "true"},
						User:    "",
						Cwd:     "",
					},
				},
				Returns: handlers.CommandHandlerCommandReturns{Err: nil},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "success(unsupported): float kwargs",
			data:      `test(name="echo", args=[1.111])`,
			mock:      handlers.CommandHandlerCommandExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "success: no args",
			data: `test("echo")`,
			mock: handlers.CommandHandlerCommandExpectation{
				Args: handlers.CommandHandlerCommandArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CommandParams{
						CmdName: "echo",
						CmdArgs: []string{},
						User:    "",
						Cwd:     "",
					},
				},
				Returns: handlers.CommandHandlerCommandReturns{Err: nil},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with user and cwd",
			data: `test("echo", [], user="testuser", cwd="/home/testuser")`,
			mock: handlers.CommandHandlerCommandExpectation{
				Args: handlers.CommandHandlerCommandArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CommandParams{
						CmdName: "echo",
						CmdArgs: []string{},
						User:    "testuser",
						Cwd:     "/home/testuser",
					},
				},
				Returns: handlers.CommandHandlerCommandReturns{Err: nil},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("echo", [], "testuser", "/home/testuser", "too many")`,
			mock:      handlers.CommandHandlerCommandExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: execute command failed",
			data: `test("echo")`,
			mock: handlers.CommandHandlerCommandExpectation{
				Args: handlers.CommandHandlerCommandArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CommandParams{
						CmdName: "echo",
						CmdArgs: []string{},
						User:    "",
						Cwd:     "",
					},
				},
				Returns: handlers.CommandHandlerCommandReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(handlers.MockCommandHandler)
			handler.ApplyCommandExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Command(handler))
			tt.errAssert(t, err)
		})
	}
}
