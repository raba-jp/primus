package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/fish/handlers"
	"github.com/raba-jp/primus/pkg/operations/fish/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/fish/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestSetVariable(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.SetVariableHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="GOPATH", value="$HOME/go", scope="universal", export=True)`,
			mock: mocks.SetVariableHandlerRunExpectation{
				Args: mocks.SetVariableHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.UniversalScope,
						Export: true,
					},
				},
				Returns: mocks.SetVariableHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: args",
			data: `test("GOPATH", "$HOME/go", "universal", True)`,
			mock: mocks.SetVariableHandlerRunExpectation{
				Args: mocks.SetVariableHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.UniversalScope,
						Export: true,
					},
				},
				Returns: mocks.SetVariableHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: global scope",
			data: `test("GOPATH", "$HOME/go", "global", True)`,
			mock: mocks.SetVariableHandlerRunExpectation{
				Args: mocks.SetVariableHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.GlobalScope,
						Export: true,
					},
				},
				Returns: mocks.SetVariableHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: local scope",
			data: `test("GOPATH", "$HOME/go", "local", True)`,
			mock: mocks.SetVariableHandlerRunExpectation{
				Args: mocks.SetVariableHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.LocalScope,
						Export: true,
					},
				},
				Returns: mocks.SetVariableHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: unexpected scope",
			data:      `test(name="GOPATH", value="$HOME/go", scope="dummy", export=True)`,
			mock:      mocks.SetVariableHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name:      "error: too many arguments",
			data:      `test("GOPATH", "$HOME/go", "universal", True, "too many")`,
			mock:      mocks.SetVariableHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: return handler error",
			data: `test("GOPATH", "$HOME/go", "universal", True)`,
			mock: mocks.SetVariableHandlerRunExpectation{
				Args: mocks.SetVariableHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.UniversalScope,
						Export: true,
					},
				},
				Returns: mocks.SetVariableHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setVariable := new(mocks.SetVariableHandler)
			setVariable.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.SetVariable(setVariable))
			tt.errAssert(t, err)
		})
	}
}
