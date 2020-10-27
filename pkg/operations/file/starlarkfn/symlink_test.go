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

func TestSymlink(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.SymlinkHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: mocks.SymlinkHandlerRunExpectation{
				Args: mocks.SymlinkHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SymlinkParams{
						Src:  "/sym/src.txt",
						Dest: "/sym/dest.txt",
					},
				},
				Returns: mocks.SymlinkHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("/sym/src.txt", "/sym/dest.txt", "too many")`,
			mock:      mocks.SymlinkHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: create symlink failed ",
			data: `test("/sym/src.txt", "/sym/dest.txt")`,
			mock: mocks.SymlinkHandlerRunExpectation{
				Args: mocks.SymlinkHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.SymlinkParams{
						Src:  "/sym/src.txt",
						Dest: "/sym/dest.txt",
					},
				},
				Returns: mocks.SymlinkHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.SymlinkHandler)
			handler.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Symlink(handler))
			tt.errAssert(t, err)
		})
	}
}
