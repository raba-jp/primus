package starlarkfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/file/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/file/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockCopyHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: func(m *mock_handlers.MockCopyHandler) {
				m.EXPECT().Copy(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.CopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o777,
						Cwd:        "/sym",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: with permission",
			data: `test("/sym/src.txt", "/sym/dest.txt", 0o644)`,
			mock: func(m *mock_handlers.MockCopyHandler) {
				m.EXPECT().Copy(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.CopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o644,
						Cwd:        "/sym",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("src.txt", "dest.txt", 0o644, "too many")`,
			mock:   func(m *mock_handlers.MockCopyHandler) {},
			hasErr: true,
		},
		{
			name: "error: file copy failed",
			data: `test("src.txt", "dest.txt", 0o644, )`,
			mock: func(m *mock_handlers.MockCopyHandler) {
				m.EXPECT().Copy(
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

			m := mock_handlers.NewMockCopyHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Copy(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
