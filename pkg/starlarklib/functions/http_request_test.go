package functions_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_executor "github.com/raba-jp/primus/pkg/executor/mock"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func TestHttpRequest(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_executor.MockExecutor)
		hasErr bool
	}{
		{
			name: "success",
			data: `http_request(url="https://example.com/", path="/sym/test.txt")`,
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().HTTPRequest(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `http_request("https://example.com/", "/sym/test.txt", "too many")`,
			mock:   func(m *mock_executor.MockExecutor) {},
			hasErr: true,
		},
		{
			name: "error: http request failed",
			data: `http_request("https://example.com/", "/sym/test.txt")`,
			mock: func(m *mock_executor.MockExecutor) {
				m.EXPECT().HTTPRequest(gomock.Any(), gomock.Any()).Return(true, xerrors.New("dummy"))
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
				"http_request": starlark.NewBuiltin("http_request", functions.HTTPRequest(m)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
