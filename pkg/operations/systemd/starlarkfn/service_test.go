package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/systemd/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/systemd/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestEnableService(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.EnableServiceHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="dummy.service")`,
			mock: mocks.EnableServiceHandlerRunExpectation{
				Args: mocks.EnableServiceHandlerRunArgs{
					CtxAnything: true,
					Name:        "dummy.service",
				},
				Returns: mocks.EnableServiceHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("dummy.service", "too many")`,
			mock:      mocks.EnableServiceHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to service enable",
			data: `test(name="dummy.service")`,
			mock: mocks.EnableServiceHandlerRunExpectation{
				Args: mocks.EnableServiceHandlerRunArgs{
					CtxAnything: true,
					Name:        "dummy.service",
				},
				Returns: mocks.EnableServiceHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enableService := new(mocks.EnableServiceHandler)
			enableService.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.EnableService(enableService))
			tt.errAssert(t, err)
		})
	}
}

func TestStartService(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.StartServiceHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="dummy.service")`,
			mock: mocks.StartServiceHandlerRunExpectation{
				Args: mocks.StartServiceHandlerRunArgs{
					CtxAnything: true,
					Name:        "dummy.service",
				},
				Returns: mocks.StartServiceHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("dummy.service", "too many")`,
			mock:      mocks.StartServiceHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to service start",
			data: `test(name="dummy.service")`,
			mock: mocks.StartServiceHandlerRunExpectation{
				Args: mocks.StartServiceHandlerRunArgs{
					CtxAnything: true,
					Name:        "dummy.service",
				},
				Returns: mocks.StartServiceHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startService := new(mocks.StartServiceHandler)
			startService.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.StartService(startService))
			tt.errAssert(t, err)
		})
	}
}
