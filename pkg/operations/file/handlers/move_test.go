package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	"github.com/spf13/afero"
)

func TestNewMove(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() afero.Fs
		params    *handlers.MoveParams
		contents  string
		errAssert assert.ErrorAssertionFunc
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
			contents:  "test",
			errAssert: assert.NoError,
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
			contents:  "test",
			errAssert: assert.NoError,
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
			contents:  "test",
			errAssert: assert.NoError,
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
			contents:  "test",
			errAssert: assert.NoError,
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
			contents:  "",
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setup()

			move := handlers.NewMove(fs)
			err := move.Run(context.Background(), tt.params)
			tt.errAssert(t, err)

			data, _ := afero.ReadFile(fs, tt.params.Dest)
			assert.Equal(t, tt.contents, string(data))

			_, err = fs.Stat(tt.params.Src)
			assert.NotNil(t, err)
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

			move := handlers.NewMove(nil)
			ctx := ctxlib.SetDryRun(context.Background(), true)
			err := move.Run(ctx, &handlers.MoveParams{
				Src:  tt.src,
				Dest: tt.dest,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
