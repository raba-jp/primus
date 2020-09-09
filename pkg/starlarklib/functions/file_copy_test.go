package functions_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/executor"
	mock_executor "github.com/raba-jp/primus/pkg/executor/mock"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func TestFileCopy(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		filename string
		mock     func(*mock_executor.MockExecutor)
		hasErr   bool
	}{
		{
			name:     "success",
			expr:     `file_copy(src="/sym/src.txt", dest="/sym/dest.txt")`,
			filename: "test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Eq(&executor.FileCopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o777,
					}),
				).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path current path",
			expr:     `file_copy("src.txt", "dest.txt")`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Eq(&executor.FileCopyParams{
						Src:        "/sym/test/src.txt",
						Dest:       "/sym/test/dest.txt",
						Permission: 0o777,
					}),
				).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path child dir",
			expr:     `file_copy("test2/src.txt", "test2/dest.txt")`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Eq(&executor.FileCopyParams{
						Src:        "/sym/test/test2/src.txt",
						Dest:       "/sym/test/test2/dest.txt",
						Permission: 0o777,
					}),
				).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path parent dir",
			expr:     `file_copy("../src.txt", "../dest.txt")`,
			filename: "/sym/test/test2/test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Eq(&executor.FileCopyParams{
						Src:        "/sym/test/src.txt",
						Dest:       "/sym/test/dest.txt",
						Permission: 0o777,
					}),
				).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:     "success: with permission",
			expr:     `file_copy("/sym/src.txt", "/sym/dest.txt", 0o644)`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Eq(&executor.FileCopyParams{
						Src:        "/sym/src.txt",
						Dest:       "/sym/dest.txt",
						Permission: 0o644,
					}),
				).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:     "error: too many arguments",
			expr:     `file_copy("src.txt", "dest.txt", 0o644, "too many")`,
			filename: "/sym/test/test.star",
			mock:     func(m *mock_executor.MockExecutor) {},
			hasErr:   true,
		},
		{
			name:     "error: file copy failed",
			expr:     `file_copy("src.txt", "dest.txt", 0o644, )`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileCopy(
					gomock.Any(),
					gomock.Any(),
				).Return(true, xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_executor.NewMockExecutor(ctrl)
			tt.mock(m)

			predeclared := starlark.StringDict{
				"file_copy": starlark.NewBuiltin("file_copy", functions.FileCopy(m)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, tt.filename, tt.expr, predeclared)
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
