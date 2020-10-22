package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/systemd/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/systemd/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestEnableService(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      func(*mock_handlers.MockEnableServiceHandler)
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="dummy.service")`,
			mock: func(m *mock_handlers.MockEnableServiceHandler) {
				m.EXPECT().EnableService(gomock.Any(), gomock.Any(), gomock.Eq("dummy.service")).Return(nil)
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("dummy.service", "too many")`,
			mock:      func(m *mock_handlers.MockEnableServiceHandler) {},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to service enable",
			data: `test(name="dummy.service")`,
			mock: func(m *mock_handlers.MockEnableServiceHandler) {
				m.EXPECT().EnableService(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockEnableServiceHandler(ctrl)

			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.EnableService(m))
			tt.errAssert(t, err)
		})
	}
}

func TestStartService(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      func(*mock_handlers.MockStartServiceHandler)
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="dummy.service")`,
			mock: func(m *mock_handlers.MockStartServiceHandler) {
				m.EXPECT().StartService(gomock.Any(), gomock.Any(), gomock.Eq("dummy.service")).Return(nil)
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("dummy.service", "too many")`,
			mock:      func(m *mock_handlers.MockStartServiceHandler) {},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to service start",
			data: `test(name="dummy.service")`,
			mock: func(m *mock_handlers.MockStartServiceHandler) {
				m.EXPECT().StartService(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockStartServiceHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.StartService(m))
			tt.errAssert(t, err)
		})
	}
}
