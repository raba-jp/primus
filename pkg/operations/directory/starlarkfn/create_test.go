package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/directory/handlers"
	"github.com/raba-jp/primus/pkg/operations/directory/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/directory/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestCreateDirectory(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.CreateHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(path="/sym/test", permission=0o777)`,
			mock: mocks.CreateHandlerRunExpectation{
				Args: mocks.CreateHandlerRunArgs{
					CtxAnything: true,
					P: &handlers.CreateParams{
						Path:       "/sym/test",
						Permission: 0o777,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CreateHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: relative path",
			data: `test(path="test", permission=0o777)`,
			mock: mocks.CreateHandlerRunExpectation{
				Args: mocks.CreateHandlerRunArgs{
					CtxAnything: true,
					P: &handlers.CreateParams{
						Path:       "test",
						Permission: 0o777,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CreateHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: without permission",
			data: `test(path="/sym/test")`,
			mock: mocks.CreateHandlerRunExpectation{
				Args: mocks.CreateHandlerRunArgs{
					CtxAnything: true,
					P: &handlers.CreateParams{
						Path:       "/sym/test",
						Permission: 0o644,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CreateHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("/sym/test", 0o644, "too many")`,
			mock:      mocks.CreateHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to create directory",
			data: `test("/sym/test", 0o644)`,
			mock: mocks.CreateHandlerRunExpectation{
				Args: mocks.CreateHandlerRunArgs{
					CtxAnything: true,
					P: &handlers.CreateParams{
						Path:       "/sym/test",
						Permission: 0o644,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CreateHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			create := new(mocks.CreateHandler)
			create.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Create(create))
			tt.errAssert(t, err)
		})
	}
}
