package functions_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_backend "github.com/raba-jp/primus/pkg/internal/backend/mock"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func TestPackage(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_backend.MockBackend)
		hasErr bool
	}{
		{
			name: "success",
			data: `package(name="base-devel")`,
			mock: func(m *mock_backend.MockBackend) {
				m.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(false)
				m.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `package("base-devel", "too many")`,
			mock:   func(m *mock_backend.MockBackend) {},
			hasErr: true,
		},
		{
			name: "error: package install failed",
			data: `package(name="base-devel")`,
			mock: func(m *mock_backend.MockBackend) {
				m.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(false)
				m.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_backend.NewMockBackend(ctrl)
			tt.mock(m)

			predeclared := starlark.StringDict{
				"package": starlark.NewBuiltin("package", functions.Package(m)),
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
