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

func TestSetPath(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.SetPathHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(values=["$GOPATH/bin", "$HOME/.bin"])`,
			mock: mocks.SetPathHandlerRunExpectation{
				Args: mocks.SetPathHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					},
				},
				Returns: mocks.SetPathHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: args",
			data: `test(["$GOPATH/bin", "$HOME/.bin"])`,
			mock: mocks.SetPathHandlerRunExpectation{
				Args: mocks.SetPathHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					},
				},
				Returns: mocks.SetPathHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: include int and bool",
			data: `test(["$GOPATH/bin", 1, True, "$HOME/.bin"])`,
			mock: mocks.SetPathHandlerRunExpectation{
				Args: mocks.SetPathHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					},
				},
				Returns: mocks.SetPathHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test(["$GOPATH/bin", "$HOME/.bin"], "too many")`,
			mock:      mocks.SetPathHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: return handler error",
			data: `test(["$GOPATH/bin", "$HOME/.bin"])`,
			mock: mocks.SetPathHandlerRunExpectation{
				Args: mocks.SetPathHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					},
				},
				Returns: mocks.SetPathHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setPath := new(mocks.SetPathHandler)
			setPath.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.SetPath(setPath))
			tt.errAssert(t, err)
		})
	}
}
