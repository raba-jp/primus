package starlarkfn_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/operations/arch/handlers"
	"github.com/raba-jp/primus/pkg/operations/arch/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/arch/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestUninstall(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   mocks.UninstallHandlerRunExpectation
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: mocks.UninstallHandlerRunExpectation{
				Args: mocks.UninstallHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.UninstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.UninstallHandlerRunReturns{
					Err: nil,
				},
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "yay", "too many")`,
			mock:   mocks.UninstallHandlerRunExpectation{},
			hasErr: true,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: mocks.UninstallHandlerRunExpectation{
				Args: mocks.UninstallHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.UninstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.UninstallHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.UninstallHandler)
			handler.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Uninstall(handler))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
