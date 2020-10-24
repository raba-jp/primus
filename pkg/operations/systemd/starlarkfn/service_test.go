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
		mock      mocks.EnableServiceHandlerEnableServiceExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="dummy.service")`,
			mock: mocks.EnableServiceHandlerEnableServiceExpectation{
				Args: mocks.EnableServiceHandlerEnableServiceArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					Name:           "dummy.service",
				},
				Returns: mocks.EnableServiceHandlerEnableServiceReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("dummy.service", "too many")`,
			mock:      mocks.EnableServiceHandlerEnableServiceExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to service enable",
			data: `test(name="dummy.service")`,
			mock: mocks.EnableServiceHandlerEnableServiceExpectation{
				Args: mocks.EnableServiceHandlerEnableServiceArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					Name:           "dummy.service",
				},
				Returns: mocks.EnableServiceHandlerEnableServiceReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.EnableServiceHandler)
			handler.ApplyEnableServiceExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.EnableService(handler))
			tt.errAssert(t, err)
		})
	}
}

func TestStartService(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.StartServiceHandlerStartServiceExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="dummy.service")`,
			mock: mocks.StartServiceHandlerStartServiceExpectation{
				Args: mocks.StartServiceHandlerStartServiceArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					Name:           "dummy.service",
				},
				Returns: mocks.StartServiceHandlerStartServiceReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("dummy.service", "too many")`,
			mock:      mocks.StartServiceHandlerStartServiceExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to service start",
			data: `test(name="dummy.service")`,
			mock: mocks.StartServiceHandlerStartServiceExpectation{
				Args: mocks.StartServiceHandlerStartServiceArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					Name:           "dummy.service",
				},
				Returns: mocks.StartServiceHandlerStartServiceReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.StartServiceHandler)
			handler.ApplyStartServiceExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.StartService(handler))
			tt.errAssert(t, err)
		})
	}
}
