package apply_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/spf13/afero"
)

func TestFileMove(t *testing.T) {
	tests := []struct {
		src      string
		dest     string
		contents string
	}{
		{
			src:      "/sym/src.txt",
			dest:     "/sym/dest.txt",
			contents: "test",
		},
	}

	for _, tt := range tests {
		defer Reset()

		_ = afero.WriteFile(fs, tt.src, []byte(tt.contents), 0777)

		exc := NewTestExecutor()
		ret, err := exc.FileMove(context.Background(), &executor.FileMoveParams{
			Src:  tt.src,
			Dest: tt.dest,
		})
		if err != nil {
			t.Fatalf("%v", err)
		}
		if !ret {
			t.Fatalf("Failed to file copy: %s %s", tt.src, tt.dest)
		}

		data, err := afero.ReadFile(fs, tt.dest)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if diff := cmp.Diff(tt.contents, string(data)); diff != "" {
			t.Fatalf("%s", diff)
		}
		if _, err := fs.Stat(tt.src); err == nil {
			t.Fatal("src file exists")
		}
	}
}
