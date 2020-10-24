package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/git/handlers"
	"github.com/raba-jp/primus/pkg/operations/git/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/git/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestClone(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.CloneHandlerCloneExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success:",
			data: `test(url="https://example.com", path="/sym", branch="main")`,
			mock: mocks.CloneHandlerCloneExpectation{
				Args: mocks.CloneHandlerCloneArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CloneParams{
						URL:    "https://example.com",
						Path:   "/sym",
						Branch: "main",
					},
				},
				Returns: mocks.CloneHandlerCloneReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: failed to git clone",
			data: `test("https://example.com", "/sym", "main")`,
			mock: mocks.CloneHandlerCloneExpectation{
				Args: mocks.CloneHandlerCloneArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CloneParams{
						URL:    "https://example.com",
						Path:   "/sym",
						Branch: "main",
					},
				},
				Returns: mocks.CloneHandlerCloneReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
		{
			name:      "error: too many arguments",
			data:      `test("https://example.com", "/sym", "main", "too many")`,
			mock:      mocks.CloneHandlerCloneExpectation{},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.CloneHandler)
			handler.ApplyCloneExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Clone(handler))
			tt.errAssert(t, err)
		})
	}
}
