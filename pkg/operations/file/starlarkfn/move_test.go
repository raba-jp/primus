package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	"github.com/raba-jp/primus/pkg/operations/file/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func TestFileMove(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      handlers.MoveHandlerMoveExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: handlers.MoveHandlerMoveExpectation{
				Args: handlers.MoveHandlerMoveArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.MoveParams{
						Src:  "/sym/src.txt",
						Dest: "/sym/dest.txt",
						Cwd:  "/sym",
					},
				},
				Returns: handlers.MoveHandlerMoveReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("src.txt", "dest.txt", "too many")`,
			mock:      handlers.MoveHandlerMoveExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: file move failed",
			data: `test("src.txt", "dest.txt")`,
			mock: handlers.MoveHandlerMoveExpectation{
				Args: handlers.MoveHandlerMoveArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.MoveParams{
						Src:  "src.txt",
						Dest: "dest.txt",
						Cwd:  "/sym",
					},
				},
				Returns: handlers.MoveHandlerMoveReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(handlers.MockMoveHandler)
			handler.ApplyMoveExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Move(handler))
			tt.errAssert(t, err)
		})
	}
}
