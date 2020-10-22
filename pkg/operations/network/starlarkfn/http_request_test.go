package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/network/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/network/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestHttpRequest(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      func(*mock_handlers.MockHTTPRequestHandler)
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(url="https://example.com/", path="/sym/test.txt")`,
			mock: func(m *mock_handlers.MockHTTPRequestHandler) {
				m.EXPECT().HTTPRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("https://example.com/", "/sym/test.txt", "too many")`,
			mock:      func(m *mock_handlers.MockHTTPRequestHandler) {},
			errAssert: assert.Error,
		},
		{
			name: "error: http request failed",
			data: `test("https://example.com/", "/sym/test.txt")`,
			mock: func(m *mock_handlers.MockHTTPRequestHandler) {
				m.EXPECT().HTTPRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockHTTPRequestHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.HTTPRequest(m))
			tt.errAssert(t, err)
		})
	}
}
