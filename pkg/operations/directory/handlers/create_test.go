package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/directory/handlers"
	"github.com/spf13/afero"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() afero.Fs
		params *handlers.CreateParams
		got    string
		hasErr bool
	}{
		{
			name: "success",
			setup: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &handlers.CreateParams{
				Path:       "/sym/test",
				Permission: 0o644,
			},
			got:    "/sym/test",
			hasErr: false,
		},
		{
			name: "success: already exists",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				fs.MkdirAll("/sym/test", 0o644)
				return fs
			},
			params: &handlers.CreateParams{
				Path:       "/sym/test",
				Permission: 0o644,
			},
			got:    "/sym/test",
			hasErr: false,
		},
		{
			name: "success: relative path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				fs.MkdirAll("/sym", 0o644)
				return fs
			},
			params: &handlers.CreateParams{
				Path:       "test",
				Permission: 0o644,
				Cwd:        "/sym",
			},
			got:    "/sym/test",
			hasErr: false,
		},
		{
			name: "success: relative path. child directory",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				fs.MkdirAll("/sym/test", 0o644)
				return fs
			},
			params: &handlers.CreateParams{
				Path:       "./test2/test3",
				Permission: 0o644,
				Cwd:        "/sym/test",
			},
			got:    "/sym/test/test2/test3",
			hasErr: false,
		},
		{
			name: "success: relative path. parent directory",
			setup: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &handlers.CreateParams{
				Path:       "../test2/test3",
				Permission: 0o644,
				Cwd:        "/sym/test",
			},
			got:    "/sym/test2/test3",
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setup()
			handler := handlers.New(fs)
			err := handler.Create(context.Background(), false, tt.params)
			if !tt.hasErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			stat, err := fs.Stat(tt.got)
			if err != nil {
				t.Fatalf("create directory failed: %v", err)
			}
			if !stat.IsDir() {
				t.Errorf("%s is not directory", tt.got)
			}
		})
	}
}

func TestNew__DryRun(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		params *handlers.CreateParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.CreateParams{
				Path:       "/sym/test",
				Permission: 0o644,
			},
			want: "mkdir -p /sym/test\nchmod 644 /sym/test\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.New(nil)
			err := handler.Create(context.Background(), true, tt.params)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
