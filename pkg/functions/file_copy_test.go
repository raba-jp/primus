package functions_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_executor "github.com/raba-jp/primus/pkg/executor/mock"
	"github.com/raba-jp/primus/pkg/functions"
	"go.starlark.net/starlark"
)

func TestFileCopy(t *testing.T) {
	tests := []string{
		`file_copy(src="/sym/src.txt", dest="/sym/dest.txt")`,
		`file_copy("/sym/src.txt", "/sym/dest.txt")`,
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_executor.NewMockExecutor(ctrl)
			m.EXPECT().FileCopy(gomock.Any(), gomock.Any()).Return(true, nil)

			predeclared := starlark.StringDict{
				"file_copy": starlark.NewBuiltin("file_copy", functions.FileCopy(context.Background(), m)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt, predeclared)
			if err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
