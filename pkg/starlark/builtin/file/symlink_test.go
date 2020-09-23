package file_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/handlers/mock"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
	"github.com/raba-jp/primus/pkg/starlark/builtin/file"
	"golang.org/x/xerrors"
)

func TestSymlink(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockSymlinkHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: func(m *mock_handlers.MockSymlinkHandler) {
				m.EXPECT().Symlink(
					gomock.Any(),
					gomock.Any(),
					&handlers.SymlinkParams{
						Src:  "/sym/src.txt",
						Dest: "/sym/dest.txt",
					},
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("/sym/src.txt", "/sys/dest.txt", "too many")`,
			mock:   func(m *mock_handlers.MockSymlinkHandler) {},
			hasErr: true,
		},
		{
			name: "error: create symlink failed ",
			data: `test("/sym/src.txt", "/sys/dest.txt")`,
			mock: func(m *mock_handlers.MockSymlinkHandler) {
				m.EXPECT().Symlink(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockSymlinkHandler(ctrl)
			tt.mock(m)

			_, err := builtin.ExecForTest("test", tt.data, file.Symlink(m))
			if !tt.hasErr && err != nil {
				t.Fatal(err)
			}
		})
	}
}
