package starlarkfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/operations/vscode/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/vscode/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/vscode/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestInstallExtension(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockInstallExtensionHandler)
		hasErr bool
	}{
		{
			name: "success: without version",
			data: `test(name="test")`,
			mock: func(m *mock_handlers.MockInstallExtensionHandler) {
				m.EXPECT().InstallExtension(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.InstallExtensionParams{
						Name: "test",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: with version",
			data: `test(name="test", version="1.0")`,
			mock: func(m *mock_handlers.MockInstallExtensionHandler) {
				m.EXPECT().InstallExtension(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.InstallExtensionParams{
						Name:    "test",
						Version: "1.0",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("test", "1.0", "too many")`,
			mock:   func(m *mock_handlers.MockInstallExtensionHandler) {},
			hasErr: true,
		},
		{
			name: "error: return handler error",
			data: `test("test")`,
			mock: func(m *mock_handlers.MockInstallExtensionHandler) {
				m.EXPECT().InstallExtension(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockInstallExtensionHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.InstallExtension(m))
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
