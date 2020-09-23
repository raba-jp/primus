package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	"github.com/spf13/afero"
)

func TestNewMove(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() afero.Fs
		params   *handlers.MoveParams
		contents string
		hasErr   bool
	}{
		{
			name: "success",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0777)
				return fs
			},
			params: &handlers.MoveParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "success: relative path current path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0777)
				return fs
			},
			params: &handlers.MoveParams{
				Src:  "src.txt",
				Dest: "dest.txt",
				Cwd:  "/sym",
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "success: relative path child path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/test/src.txt", []byte("test"), 0777)
				return fs
			},
			params: &handlers.MoveParams{
				Src:  "./test/src.txt",
				Dest: "./test/dest.txt",
				Cwd:  "/sym",
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "success: relative path parent path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/test2/src.txt", []byte("test"), 0777)
				return fs
			},
			params: &handlers.MoveParams{
				Src:  "../test2/src.txt",
				Dest: "../test2/dest.txt",
				Cwd:  "/sym/test",
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "error: source file not found",
			setup: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &handlers.MoveParams{
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

			handler := handlers.NewMove(fs)
			err := handler.Move(context.Background(), false, tt.params)
			if !tt.hasErr {
				if err != nil {
					t.Fatalf("%v", err)
				}

				data, err := afero.ReadFile(fs, tt.params.Dest)
				if err != nil {
					t.Fatalf("dest file read failed: %s: %v", tt.params.Dest, err)
				}
				if diff := cmp.Diff(tt.contents, string(data)); diff != "" {
					t.Fatal(diff)
				}
				if _, err := fs.Stat(tt.params.Src); err == nil {
					t.Fatal("src file exists")
				}
			}
		})
	}
}

func TestNewMove__DryRun(t *testing.T) {
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
			want: "mv /sym/src.txt /sym/dest.txt\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.NewMove(nil)
			err := handler.Move(context.Background(), true, &handlers.MoveParams{
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
