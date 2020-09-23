package file_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/handlers/mock"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/raba-jp/primus/pkg/starlark/builtin/file"
	"golang.org/x/xerrors"
)

func TestFileMove(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockFileMoveHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
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
		// {
		// 	name: "success: relative path current path",
		// 	data: `test("src.txt", "dest.txt")`,
		// 	mock: func(m *mock_handlers.MockFileMoveHandler) {
		// 		m.EXPECT().FileMove(
		// 			gomock.Any(),
		// 			gomock.Any(),
		// 			gomock.Eq(&handlers.FileMoveParams{
		// 				Src:  "/sym/test/src.txt",
		// 				Dest: "/sym/test/dest.txt",
		// 			}),
		// 		).Return(nil)
		// 	},
		// 	hasErr: false,
		// },
		// {
		// 	name: "success: relative path child dir",
		// 	data: `test("test2/src.txt", "test2/dest.txt")`,
		// 	mock: func(m *mock_handlers.MockFileMoveHandler) {
		// 		m.EXPECT().FileMove(
		// 			gomock.Any(),
		// 			gomock.Any(),
		// 			gomock.Eq(&handlers.FileMoveParams{
		// 				Src:  "/sym/test/test2/src.txt",
		// 				Dest: "/sym/test/test2/dest.txt",
		// 			}),
		// 		).Return(nil)
		// 	},
		// 	hasErr: false,
		// },
		// {
		// 	name: "success: relative path parent dir",
		// 	data: `test("../src.txt", "../dest.txt")`,
		// 	mock: func(m *mock_handlers.MockFileMoveHandler) {
		// 		m.EXPECT().FileMove(
		// 			gomock.Any(),
		// 			gomock.Any(),
		// 			gomock.Eq(&handlers.FileMoveParams{
		// 				Src:  "/sym/test/src.txt",
		// 				Dest: "/sym/test/dest.txt",
		// 			}),
		// 		).Return(nil)
		// 	},
		// 	hasErr: false,
		// },
		{
			name:   "error: too many arguments",
			data:   `test("src.txt", "dest.txt", "too many")`,
			mock:   func(m *mock_handlers.MockFileMoveHandler) {},
			hasErr: true,
		},
		{
			name: "error: file move failed",
			data: `test("src.txt", "dest.txt")`,
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

			_, err := starlark.ExecForTest("test", tt.data, file.Move(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
