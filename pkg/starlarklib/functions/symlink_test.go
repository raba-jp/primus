package functions_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_backend "github.com/raba-jp/primus/pkg/internal/backend/mock"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func TestSymlink(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_backend.MockBackend)
		hasErr bool
	}{
		{
			name: "success",
			data: `symlink(src="/sym/src.txt", dest="/sys/dest.txt")`,
			mock: func(m *mock_backend.MockBackend) {
				m.EXPECT().Symlink(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `symlink("/sym/src.txt", "/sys/dest.txt", "too many")`,
			mock:   func(m *mock_backend.MockBackend) {},
			hasErr: true,
		},
		{
			name: "error: create symlink failed ",
			data: `symlink("/sym/src.txt", "/sys/dest.txt")`,
			mock: func(m *mock_backend.MockBackend) {
				m.EXPECT().Symlink(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
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
				"symlink": starlark.NewBuiltin("symlink", functions.Symlink(m)),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if !tt.hasErr && err != nil {
				t.Fatal(err)
			}
		})
	}
}
