package starlarkfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/packages/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestDarwinPkgInstall(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockDarwinPkgInstallHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="base-devel", option="option", cask=True, cmd="brew")`,
			mock: func(m *mock_handlers.MockDarwinPkgInstallHandler) {
				m.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "option", True, "brew", "too many")`,
			mock:   func(m *mock_handlers.MockDarwinPkgInstallHandler) {},
			hasErr: true,
		},
		{
			name: "error: package install failed",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockDarwinPkgInstallHandler) {
				m.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockDarwinPkgInstallHandler(ctrl)

			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.DarwinPkgInstall(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestArchPkgInstall(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockArchPkgInstallHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="base-devel", option="option", cmd="yay")`,
			mock: func(m *mock_handlers.MockArchPkgInstallHandler) {
				m.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "option", "yay", "too many")`,
			mock:   func(m *mock_handlers.MockArchPkgInstallHandler) {},
			hasErr: true,
		},
		{
			name: "error: package install failed",
			data: `test(name="base-devel")`,
			mock: func(m *mock_handlers.MockArchPkgInstallHandler) {
				m.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockArchPkgInstallHandler(ctrl)

			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.ArchPkgInstall(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
