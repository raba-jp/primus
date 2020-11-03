package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/directory/handlers"
	"github.com/spf13/afero"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() afero.Fs
		params    *handlers.CreateParams
		got       string
		errAssert assert.ErrorAssertionFunc
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
			got:       "/sym/test",
			errAssert: assert.NoError,
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
			got:       "/sym/test",
			errAssert: assert.NoError,
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
			got:       "/sym/test",
			errAssert: assert.NoError,
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
			got:       "/sym/test/test2/test3",
			errAssert: assert.NoError,
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
			got:       "/sym/test2/test3",
			errAssert: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setup()
			create := handlers.New(fs)
			err := create.Run(context.Background(), tt.params)
			tt.errAssert(t, err)

			stat, err := fs.Stat(tt.got)
			assert.NoErrorf(t, err, "create directory failed")

			isdir := !stat.IsDir()
			assert.Falsef(t, isdir, tt.got+"is not directory")
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

			create := handlers.New(nil)
			ctx := ctxlib.SetDryRun(context.Background(), true)
			err := create.Run(ctx, tt.params)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
