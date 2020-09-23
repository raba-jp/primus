package handlers_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	"github.com/spf13/afero"
)

func TestNewCopy(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() afero.Fs
		params     *handlers.CopyParams
		permission os.FileMode
		contents   string
		hasErr     bool
	}{
		{
			name: "success",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &handlers.CopyParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			permission: 0o777,
			contents:   "test",
			hasErr:     false,
		},
		{
			name: "success: set permission",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &handlers.CopyParams{
				Src:        "/sym/src.txt",
				Dest:       "/sym/dest.txt",
				Permission: 0o644,
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "success: relative path current path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &handlers.CopyParams{
				Src:        "./src.txt",
				Dest:       "./dest.txt",
				Permission: 0o777,
				Cwd:        "/sym",
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "success: relative path child path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/sym/test/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &handlers.CopyParams{
				Src:        "./test/src.txt",
				Dest:       "./test/dest.txt",
				Permission: 0o777,
				Cwd:        "/sym",
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "success: relative path parent path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				fs.MkdirAll("/sym/test", 0o777)
				afero.WriteFile(fs, "/sym/test2/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &handlers.CopyParams{
				Src:        "../test2/src.txt",
				Dest:       "../test2/dest.txt",
				Permission: 0o777,
				Cwd:        "/sym/test",
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "error: source file not found",
			setup: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &handlers.CopyParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			contents: "test",
			hasErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setup()
			handler := handlers.NewCopy(fs)
			err := handler.Copy(context.Background(), false, tt.params)
			if !tt.hasErr {
				if err != nil {
					t.Fatalf("%v", err)
				}

				data, err := afero.ReadFile(fs, tt.params.Dest)
				if err != nil {
					t.Fatalf("Failed to read file: %s: %v", tt.params.Dest, err)
				}
				if diff := cmp.Diff(tt.contents, string(data)); diff != "" {
					t.Fatal(diff)
				}
				stat, err := fs.Stat(tt.params.Dest)
				if err != nil {
					t.Fatalf("%v", err)
				}
				if stat.Mode() != tt.params.Permission {
					t.Fatalf("Set permission failed: %s", tt.params.Dest)
				}
			}
		})
	}
}

func TestBaseBackend_FileCopy__DryRun(t *testing.T) {
	tests := []struct {
		name string
		src  string
		dest string
		want string
	}{
		{
			name: "success",
			src:  "/sym/src.txt",
			dest: "/sym/dest.txt",
			want: "cp /sym/src.txt /sym/dest.txt\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.NewCopy(nil)
			err := handler.Copy(context.Background(), true, &handlers.CopyParams{
				Src:  tt.src,
				Dest: tt.dest,
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
