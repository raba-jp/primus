package starlarkfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/packages/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func TestDarwinPkgCheckInstall(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockDarwinPkgCheckInstallHandler)
		hasErr bool
	}{
		{
			name: "success: true",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockDarwinPkgCheckInstallHandler) {
				m.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(true)
			},
			hasErr: false,
		},
		{
			name: "success: false",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockDarwinPkgCheckInstallHandler) {
				m.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(false)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "too many")`,
			mock:   func(m *mock_handlers.MockDarwinPkgCheckInstallHandler) {},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockDarwinPkgCheckInstallHandler(ctrl)

			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.DarwinPkgCheckInstall(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestArchPkgCheckInstall(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockArchPkgCheckInstallHandler)
		hasErr bool
	}{
		{
			name: "success: true",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockArchPkgCheckInstallHandler) {
				m.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(true)
			},
			hasErr: false,
		},
		{
			name: "success: false",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockArchPkgCheckInstallHandler) {
				m.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(false)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "too many")`,
			mock:   func(m *mock_handlers.MockArchPkgCheckInstallHandler) {},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockArchPkgCheckInstallHandler(ctrl)

			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.ArchPkgCheckInstall(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
