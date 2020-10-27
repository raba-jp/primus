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
		mock      mocks.CopyHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: mocks.CopyHandlerRunExpectation{
				Args: mocks.CopyHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o777,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CopyHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with permission",
			data: `test("/sym/src.txt", "/sym/dest.txt", 0o644)`,
			mock: mocks.CopyHandlerRunExpectation{
				Args: mocks.CopyHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o644,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CopyHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("src.txt", "dest.txt", 0o644, "too many")`,
			mock:      mocks.CopyHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: file copy failed",
			data: `test("src.txt", "dest.txt", 0o644, )`,
			mock: mocks.CopyHandlerRunExpectation{
				Args: mocks.CopyHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.CopyParams{
						Src:        "src.txt",
						Dest:       "dest.txt",
						Permission: 0o644,
						Cwd:        "/sym",
					},
				},
				Returns: mocks.CopyHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copy := new(mocks.CopyHandler)
			copy.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Copy(copy))
			tt.errAssert(t, err)
		})
	}
}
