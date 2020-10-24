package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/network/handlers"
	"github.com/raba-jp/primus/pkg/operations/network/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestHttpRequest(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      handlers.HTTPRequestHandlerHTTPRequestExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(url="https://example.com/", path="/sym/test.txt")`,
			mock: handlers.HTTPRequestHandlerHTTPRequestExpectation{
				Args: handlers.HTTPRequestHandlerHTTPRequestArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.HTTPRequestParams{
						URL:  "https://example.com/",
						Path: "/sym/test.txt",
					},
				},
				Returns: handlers.HTTPRequestHandlerHTTPRequestReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("https://example.com/", "/sym/test.txt", "too many")`,
			mock:      handlers.HTTPRequestHandlerHTTPRequestExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: http request failed",
			data: `test("https://example.com/", "/sym/test.txt")`,
			mock: handlers.HTTPRequestHandlerHTTPRequestExpectation{
				Args: handlers.HTTPRequestHandlerHTTPRequestArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.HTTPRequestParams{
						URL:  "https://example.com/",
						Path: "/sym/test.txt",
					},
				},
				Returns: handlers.HTTPRequestHandlerHTTPRequestReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(handlers.MockHTTPRequestHandler)
			handler.ApplyHTTPRequestExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.HTTPRequest(handler))
			tt.errAssert(t, err)
		})
	}
}
