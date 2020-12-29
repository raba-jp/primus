package network_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/raba-jp/primus/pkg/functions/network"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

type MockRoundTripper struct {
	http.RoundTripper
	Fn func(req *http.Request) *http.Response
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Fn(req), nil
}

func MockHttpClient(fn func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{Fn: fn},
	}
}

func TestNewHttpRequestFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      network.HTTPRequestRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(url="https://example.com/", path="/sym/test.txt")`,
			mock: func(ctx context.Context, p *network.HTTPRequestParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("https://example.com/", "/sym/test.txt", "too many")`,
			mock: func(ctx context.Context, p *network.HTTPRequestParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: http request failed",
			data: `test("https://example.com/", "/sym/test.txt")`,
			mock: func(ctx context.Context, p *network.HTTPRequestParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, network.NewHTTPRequestFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestHTTPRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		url      string
		path     string
		contents string
		httpMock func(req *http.Request) *http.Response
	}{
		{
			name:     "success",
			url:      "https://example.com/",
			path:     "/sym/test.txt",
			contents: "test file",
			httpMock: func(req *http.Request) *http.Response {
				buf := bytes.NewBufferString("test file")
				body := ioutil.NopCloser(buf)
				return &http.Response{
					Body: body,
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			run := network.HTTPRequest(MockHttpClient(tt.httpMock), fs)
			err := run(context.Background(), &network.HTTPRequestParams{
				URL:  tt.url,
				Path: tt.path,
			})
			assert.NoError(t, err)

			data, err := afero.ReadFile(fs, tt.path)
			assert.NoError(t, err)
			assert.Equal(t, tt.contents, string(data))
		})
	}
}
