package starlarkfn_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/operations/darwin/handlers"
	"github.com/raba-jp/primus/pkg/operations/darwin/starlarkfn"

	"github.com/raba-jp/primus/pkg/operations/darwin/handlers/mocks"

	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestInstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.InstallHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="base-devel", option="option", cask=True, cmd="brew")`,
			mock: mocks.InstallHandlerRunExpectation{
				Args: mocks.InstallHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.InstallParams{
						Name:   "base-devel",
						Option: "option",
						Cask:   true,
						Cmd:    "brew",
					},
				},
				Returns: mocks.InstallHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "option", True, "brew", "too many")`,
			mock:      mocks.InstallHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: package install failed",
			data: `test(name="base-devel")`,
			mock: mocks.InstallHandlerRunExpectation{
				Args: mocks.InstallHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.InstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.InstallHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.InstallHandler)
			handler.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Install(handler))
			tt.errAssert(t, err)
		})
	}
}
