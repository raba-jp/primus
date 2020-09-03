package functions_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_executor "github.com/raba-jp/primus/pkg/executor/mock"
	"github.com/raba-jp/primus/pkg/functions"
	"go.starlark.net/starlark"
)

func TestPackage(t *testing.T) {
	tests := []string{
		`package(name="base-devel")`,
		`package("base-devel")`,
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_executor.NewMockExecutor(ctrl)
			m.EXPECT().Package(gomock.Any(), gomock.Any()).Return(true, nil)

			predeclared := starlark.StringDict{
				"package": starlark.NewBuiltin("package", functions.Package(context.Background(), m)),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt, predeclared)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
