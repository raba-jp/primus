package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	"github.com/raba-jp/primus/pkg/operations/file/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/file/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.CopyHandlerCopyExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: mocks.CopyHandlerCopyExpectation{
				Args: mocks.CopyHandlerCopyArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o777,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CopyHandlerCopyReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with permission",
			data: `test("/sym/src.txt", "/sym/dest.txt", 0o644)`,
			mock: mocks.CopyHandlerCopyExpectation{
				Args: mocks.CopyHandlerCopyArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o644,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CopyHandlerCopyReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("src.txt", "dest.txt", 0o644, "too many")`,
			mock:      mocks.CopyHandlerCopyExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: file copy failed",
			data: `test("src.txt", "dest.txt", 0o644, )`,
			mock: mocks.CopyHandlerCopyExpectation{
				Args: mocks.CopyHandlerCopyArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CopyParams{
						Src:        "src.txt",
						Dest:       "dest.txt",
						Permission: 0o644,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CopyHandlerCopyReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.CopyHandler)
			handler.ApplyCopyExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Copy(handler))
			tt.errAssert(t, err)
		})
	}
}
