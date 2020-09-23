package handlers_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/network/handlers"
	"github.com/spf13/afero"
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

func TestNewHTTPRequest(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			handler := handlers.NewHTTPRequest(MockHttpClient(tt.httpMock), fs)
			err := handler.HTTPRequest(context.Background(), false, &handlers.HTTPRequestParams{
				URL:  tt.url,
				Path: tt.path,
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
			data, err := afero.ReadFile(fs, tt.path)
			if err != nil {
				t.Fatalf("file read failed: %s: %v", tt.path, err)
			}
			if diff := cmp.Diff(tt.contents, string(data)); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestNewHTTPRequest__DryRun(t *testing.T) {
	tests := []struct {
		name string
		url  string
		path string
		want string
	}{
		{
			name: "success",
			url:  "https://example.com",
			path: "/sym/output.txt",
			want: "curl -Lo /sym/output.txt https://example.com\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.NewHTTPRequest(nil, nil)
			err := handler.HTTPRequest(context.Background(), true, &handlers.HTTPRequestParams{
				URL:  tt.url,
				Path: tt.path,
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
