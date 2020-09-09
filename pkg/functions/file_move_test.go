package functions_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/executor"
	mock_executor "github.com/raba-jp/primus/pkg/executor/mock"
	"github.com/raba-jp/primus/pkg/functions"
	"go.starlark.net/starlark"
)

func TestFileMove(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		filename string
		matcher  gomock.Matcher
	}{
		{
			name:     "success kwargs",
			expr:     `file_move(src="/sym/src.txt", dest="/sym/dest.txt")`,
			filename: "test.star",
			matcher: gomock.Eq(&executor.FileMoveParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			}),
		},
		{
			name:     "success args",
			expr:     `file_move("/sym/src.txt", "/sym/dest.txt")`,
			filename: "test.star",
			matcher: gomock.Eq(&executor.FileMoveParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			}),
		},
		{
			name:     "success relative path current path",
			expr:     `file_move("src.txt", "dest.txt")`,
			filename: "/sym/test/test.star",
			matcher: gomock.Eq(&executor.FileMoveParams{
				Src:  "/sym/test/src.txt",
				Dest: "/sym/test/dest.txt",
			}),
		},
		{
			name:     "success relative path child dir",
			expr:     `file_move("test2/src.txt", "test2/dest.txt")`,
			filename: "/sym/test/test.star",
			matcher: gomock.Eq(&executor.FileMoveParams{
				Src:  "/sym/test/test2/src.txt",
				Dest: "/sym/test/test2/dest.txt",
			}),
		},
		{
			name:     "success relative path parent dir",
			expr:     `file_move("../src.txt", "../dest.txt")`,
			filename: "/sym/test/test2/test.star",
			matcher: gomock.Eq(&executor.FileMoveParams{
				Src:  "/sym/test/src.txt",
				Dest: "/sym/test/dest.txt",
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_executor.NewMockExecutor(ctrl)
			m.EXPECT().FileMove(gomock.Any(), tt.matcher).Return(true, nil)

			predeclared := starlark.StringDict{
				"file_move": starlark.NewBuiltin("file_move", functions.FileMove(m)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, tt.filename, tt.expr, predeclared)
			if err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
