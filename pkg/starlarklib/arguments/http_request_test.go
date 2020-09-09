package arguments_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
)

func TestNewHTTPRequestArguments(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   *arguments.HTTPRequestArguments
		hasErr bool
	}{
		{
			name: "success: kwargs",
			data: `http_request(url="https://example.com", path="/sym/path.txt")`,
			want: &arguments.HTTPRequestArguments{
				URL:  "https://example.com",
				Path: "/sym/path.txt",
			},
			hasErr: false,
		},
		{
			name: "success: args",
			data: `http_request("https://example.com", "/sym/path.txt")`,
			want: &arguments.HTTPRequestArguments{
				URL:  "https://example.com",
				Path: "/sym/path.txt",
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `http_request("https://example.com", "/sym/path.txt", "too many")`,
			want:   &arguments.HTTPRequestArguments{},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *arguments.HTTPRequestArguments

			predeclared := starlark.StringDict{
				"http_request": starlark.NewBuiltin("http_request", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var err error
					got, err = arguments.NewHTTPRequestArguments(b, args, kwargs)
					return starlark.None, err
				}),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if !tt.hasErr && err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
