package functions_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_handlers "github.com/raba-jp/primus/pkg/internal/handlers/mock"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func TestPackage(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockCheckInstallHandler, *mock_handlers.MockInstallHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `package(name="base-devel")`,
			mock: func(ch *mock_handlers.MockCheckInstallHandler, i *mock_handlers.MockInstallHandler) {
				ch.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(false)
				i.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `package("base-devel", "too many")`,
			mock:   func(ch *mock_handlers.MockCheckInstallHandler, i *mock_handlers.MockInstallHandler) {},
			hasErr: true,
		},
		{
			name: "error: package install failed",
			data: `package(name="base-devel")`,
			mock: func(ch *mock_handlers.MockCheckInstallHandler, i *mock_handlers.MockInstallHandler) {
				ch.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(false)
				i.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ch := mock_handlers.NewMockCheckInstallHandler(ctrl)
			i := mock_handlers.NewMockInstallHandler(ctrl)

			tt.mock(ch, i)

			predeclared := starlark.StringDict{
				"package": starlark.NewBuiltin("package", functions.Package(ch, i)),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if !tt.hasErr && err != nil {
				t.Error(err)
			}
		})
	}
}
