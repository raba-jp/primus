package functions_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/starlark_iac/functions"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
)

func TestFileMove(t *testing.T) {
	tests := []struct {
		data     string
		src      string
		dest     string
		contents string
	}{
		{
			data:     `file_move(src="/sym/src.txt", dest="/sym/dest.txt")`,
			src:      "/sym/src.txt",
			dest:     "/sym/dest.txt",
			contents: "file",
		},
	}

	fs := afero.NewMemMapFs()

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			_ = afero.WriteFile(fs, tt.src, []byte(tt.contents), 0777)

			predeclared := starlark.StringDict{
				"file_move": starlark.NewBuiltin("file_move", functions.FileMove(context.Background(), fs)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if err != nil {
				t.Fatalf("%v", err)
			}
			data, err := afero.ReadFile(fs, tt.dest)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.contents, string(data)); diff != "" {
				t.Fatalf("%s", diff)
			}
			if _, err := fs.Stat(tt.src); err == nil {
				t.Fatal("src file exist")
			}
		})
	}
}
