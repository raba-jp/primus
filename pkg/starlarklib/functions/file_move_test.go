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

func TestFileMove(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		filename string
		mock     func(*mock_handlers.MockFileMoveHandler)
		hasErr   bool
	}{
		{
			name:     "success",
			expr:     `move_file(src="/sym/src.txt", dest="/sym/dest.txt")`,
			filename: "test.star",
			mock: func(m *mock_handlers.MockFileMoveHandler) {
				m.EXPECT().FileMove(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FileMoveParams{
						Src:  "/sym/src.txt",
						Dest: "/sym/dest.txt",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path current path",
			expr:     `move_file("src.txt", "dest.txt")`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_handlers.MockFileMoveHandler) {
				m.EXPECT().FileMove(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FileMoveParams{
						Src:  "/sym/test/src.txt",
						Dest: "/sym/test/dest.txt",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path child dir",
			expr:     `move_file("test2/src.txt", "test2/dest.txt")`,
			filename: "/sym/test/test.star",
			mock: func(m *mock_handlers.MockFileMoveHandler) {
				m.EXPECT().FileMove(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FileMoveParams{
						Src:  "/sym/test/test2/src.txt",
						Dest: "/sym/test/test2/dest.txt",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:     "success: relative path parent dir",
			expr:     `move_file("../src.txt", "../dest.txt")`,
			filename: "/sym/test/test2/test.star",
			mock: func(m *mock_handlers.MockFileMoveHandler) {
				m.EXPECT().FileMove(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FileMoveParams{
						Src:  "/sym/test/src.txt",
						Dest: "/sym/test/dest.txt",
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:     "error: too many arguments",
			expr:     `move_file("src.txt", "dest.txt", "too many")`,
			filename: "/sym/test/test2/test.star",
			mock:     func(m *mock_handlers.MockFileMoveHandler) {},
			hasErr:   true,
		},
		{
			name:     "error: file move failed",
			expr:     `move_file("src.txt", "dest.txt")`,
			filename: "/sym/test/test2/test.star",
			mock: func(m *mock_handlers.MockFileMoveHandler) {
				m.EXPECT().FileMove(
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

			m := mock_handlers.NewMockFileMoveHandler(ctrl)
			tt.mock(m)

			predeclared := starlark.StringDict{
				"move_file": starlark.NewBuiltin("move_file", functions.FileMove(m)),
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
