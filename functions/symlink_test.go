package functions_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_executor "github.com/raba-jp/primus/executor/mock"
	"github.com/raba-jp/primus/functions"
	"go.starlark.net/starlark"
)

func TestSymlink(t *testing.T) {
	tests := []string{
		`symlink(src="/sym/src.txt", dest="/sys/dest.txt")`,
		`symlink("/sym/src.txt", "/sys/dest.txt")`,
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_executor.NewMockExecutor(ctrl)
			m.EXPECT().Symlink(gomock.Any(), gomock.Any()).Return(true, nil)

			predeclared := starlark.StringDict{
				"symlink": starlark.NewBuiltin("symlink", functions.Symlink(context.Background(), m)),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt, predeclared)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
