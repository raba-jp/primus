package starlarkfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/packages/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestDarwinPkgUninstall(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockDarwinPkgUninstallHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockDarwinPkgUninstallHandler) {
				m.EXPECT().Uninstall(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "too many")`,
			mock:   func(m *mock_handlers.MockDarwinPkgUninstallHandler) {},
			hasErr: true,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockDarwinPkgUninstallHandler) {
				m.EXPECT().Uninstall(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockDarwinPkgUninstallHandler(ctrl)

			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.DarwinPkgUninstall(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestArchPkgUninstall(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockArchPkgUninstallHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockArchPkgUninstallHandler) {
				m.EXPECT().Uninstall(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "yay", "too many")`,
			mock:   func(m *mock_handlers.MockArchPkgUninstallHandler) {},
			hasErr: true,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockArchPkgUninstallHandler) {
				m.EXPECT().Uninstall(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockArchPkgUninstallHandler(ctrl)

			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.ArchPkgUninstall(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
