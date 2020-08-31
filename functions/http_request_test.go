package functions_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_executor "github.com/raba-jp/primus/executor/mock"
	"github.com/raba-jp/primus/functions"
	"go.starlark.net/starlark"
)

func TestHttpRequest(t *testing.T) {
	tests := []string{
		`http_request(url="https://example.com/", path="/sym/test.txt")`,
		`http_request("https://example.com/", "/sym/test.txt")`,
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_executor.NewMockExecutor(ctrl)
			m.EXPECT().HTTPRequest(gomock.Any(), gomock.Any()).Return(true, nil)

			predeclared := starlark.StringDict{
				"http_request": starlark.NewBuiltin("http_request", functions.HTTPRequest(context.Background(), m)),
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
