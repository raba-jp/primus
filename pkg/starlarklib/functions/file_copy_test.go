package functions_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/handlers/mock"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func TestFileCopy(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		filename string
		mock     func(*mock_handlers.MockFileCopyHandler)
		hasErr   bool
	}{
		{
			name:     "success",
			data:     `copy_file(src="/sym/src.txt", dest="/sym/dest.txt")`,
			filename: "test.star",
			mock: func(m *mock_handlers.MockFileCopyHandler) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FileCopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o777,
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path current path",
			data:     `copy_file("src.txt", "dest.txt")`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_handlers.MockFileCopyHandler) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FileCopyParams{
						Src:        "/sym/test/src.txt",
						Dest:       "/sym/test/dest.txt",
						Permission: 0o777,
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path child dir",
			data:     `copy_file("test2/src.txt", "test2/dest.txt")`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_handlers.MockFileCopyHandler) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FileCopyParams{
						Src:        "/sym/test/test2/src.txt",
						Dest:       "/sym/test/test2/dest.txt",
						Permission: 0o777,
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path parent dir",
			data:     `copy_file("../src.txt", "../dest.txt")`,
			filename: "/sym/test/test2/test.star",
			mock: func(m *mock_handlers.MockFileCopyHandler) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FileCopyParams{
						Src:        "/sym/test/src.txt",
						Dest:       "/sym/test/dest.txt",
						Permission: 0o777,
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:     "success: with permission",
			data:     `copy_file("/sym/src.txt", "/sym/dest.txt", 0o644)`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_handlers.MockFileCopyHandler) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FileCopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o644,
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:     "error: too many arguments",
			data:     `copy_file("src.txt", "dest.txt", 0o644, "too many")`,
			filename: "/sym/test/test.star",
			mock:     func(m *mock_handlers.MockFileCopyHandler) {},
			hasErr:   true,
		},
		{
			name:     "error: file copy failed",
			data:     `copy_file("src.txt", "dest.txt", 0o644, )`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_handlers.MockFileCopyHandler) {
				m.EXPECT().FileCopy(
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

			m := mock_handlers.NewMockFileCopyHandler(ctrl)
			tt.mock(m)

			predeclared := starlark.StringDict{
				"copy_file": starlark.NewBuiltin("copy_file", functions.FileCopy(m)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, tt.filename, tt.data, predeclared)
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
