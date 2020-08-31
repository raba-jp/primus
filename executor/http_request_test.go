package executor_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/executor"
	"github.com/spf13/afero"
)

func TestHttpRequest(t *testing.T) {
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
			defer Reset()

			exc := NewTestExecutor()
			ret, err := exc.HTTPRequest(context.Background(), &executor.HTTPRequestParams{
				URL:  tt.url,
				Path: tt.path,
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
			if !ret {
				t.Fatalf("Failed to http request: %v ", tt)
			}
			data, err := afero.ReadFile(fs, tt.path)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.contents, string(data)); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
