package directory_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/handlers/mock"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/raba-jp/primus/pkg/starlark/builtin/directory"
	"golang.org/x/xerrors"
)

func TestCreateDirectory(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockCreateDirectoryHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(path="/sym/test", permission=0o777)`,
			mock: func(m *mock_handlers.MockCreateDirectoryHandler) {
				m.EXPECT().CreateDirectory(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.CreateDirectoryParams{
						Path:       "/sym/test",
						Permission: 0o777,
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: relative path",
			data: `test(path="test", permission=0o777)`,
			mock: func(m *mock_handlers.MockCreateDirectoryHandler) {
				m.EXPECT().CreateDirectory(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.CreateDirectoryParams{
						Path:       "/sym/test",
						Permission: 0o777,
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: without permission",
			data: `test(path="/sym/test")`,
			mock: func(m *mock_handlers.MockCreateDirectoryHandler) {
				m.EXPECT().CreateDirectory(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(
						&handlers.CreateDirectoryParams{
							Path:       "/sym/test",
							Permission: 0o644,
						}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("/sym/test", 0o644, "too many")`,
			mock:   func(m *mock_handlers.MockCreateDirectoryHandler) {},
			hasErr: true,
		},
		{
			name: "error: failed to create directory",
			data: `test("/sym/test", 0o644)`,
			mock: func(m *mock_handlers.MockCreateDirectoryHandler) {
				m.EXPECT().CreateDirectory(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockCreateDirectoryHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, directory.Create(m))
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
