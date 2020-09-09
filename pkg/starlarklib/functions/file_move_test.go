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

func TestFileMove(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		filename string
		mock     func(*mock_executor.MockExecutor)
		hasErr   bool
	}{
		{
			name:     "success",
			expr:     `file_move(src="/sym/src.txt", dest="/sym/dest.txt")`,
			filename: "test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileMove(
					gomock.Any(),
					gomock.Eq(&executor.FileMoveParams{
						Src:  "/sym/src.txt",
						Dest: "/sym/dest.txt",
					}),
				).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path current path",
			expr:     `file_move("src.txt", "dest.txt")`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileMove(
					gomock.Any(),
					gomock.Eq(&executor.FileMoveParams{
						Src:  "/sym/test/src.txt",
						Dest: "/sym/test/dest.txt",
					}),
				).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path child dir",
			expr:     `file_move("test2/src.txt", "test2/dest.txt")`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileMove(
					gomock.Any(),
					gomock.Eq(&executor.FileMoveParams{
						Src:  "/sym/test/test2/src.txt",
						Dest: "/sym/test/test2/dest.txt",
					}),
				).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path parent dir",
			expr:     `file_move("../src.txt", "../dest.txt")`,
			filename: "/sym/test/test2/test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileMove(
					gomock.Any(),
					gomock.Eq(&executor.FileMoveParams{
						Src:  "/sym/test/src.txt",
						Dest: "/sym/test/dest.txt",
					}),
				).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:     "error: too many arguments",
			expr:     `file_move("src.txt", "dest.txt", "too many")`,
			filename: "/sym/test/test2/test.star",
			mock:     func(m *mock_executor.MockExecutor) {},
			hasErr:   true,
		},
		{
			name:     "error: file move failed",
			expr:     `file_move("src.txt", "dest.txt")`,
			filename: "/sym/test/test2/test.star",
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().FileMove(
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
				"file_move": starlark.NewBuiltin("file_move", functions.FileMove(m)),
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
