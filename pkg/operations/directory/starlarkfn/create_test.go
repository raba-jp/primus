package starlarkfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/operations/directory/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/directory/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/directory/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestCreateDirectory(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockCreateHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(path="/sym/test", permission=0o777)`,
			mock: func(m *mock_handlers.MockCreateHandler) {
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.CreateParams{
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
			mock: func(m *mock_handlers.MockCreateHandler) {
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.CreateParams{
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
			mock: func(m *mock_handlers.MockCreateHandler) {
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(
						&handlers.CreateParams{
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
			mock:   func(m *mock_handlers.MockCreateHandler) {},
			hasErr: true,
		},
		{
			name: "error: failed to create directory",
			data: `test("/sym/test", 0o644)`,
			mock: func(m *mock_handlers.MockCreateHandler) {
				m.EXPECT().Create(
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

			m := mock_handlers.NewMockCreateHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Create(m))
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
