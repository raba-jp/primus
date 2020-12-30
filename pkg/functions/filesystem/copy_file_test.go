package filesystem_test

import (
	"context"
	"os"
	"testing"

	"github.com/raba-jp/primus/pkg/functions/filesystem"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewCopyFileFunction(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      filesystem.CopyFileRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: func(ctx context.Context, p *filesystem.CopyFileParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with permission",
			data: `test("/sym/src.txt", "/sym/dest.txt", 0o644)`,
			mock: func(ctx context.Context, p *filesystem.CopyFileParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("src.txt", "dest.txt", 0o644, "too many")`,
			mock: func(ctx context.Context, p *filesystem.CopyFileParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: file copy failed",
			data: `test("src.txt", "dest.txt", 0o644, )`,
			mock: func(ctx context.Context, p *filesystem.CopyFileParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := starlark.ExecForTest("test", tt.data, filesystem.NewCopyFileFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestCopyFile(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() afero.Fs
		params     *filesystem.CopyFileParams
		permission os.FileMode
		contents   string
		errAssert  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &filesystem.CopyFileParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			permission: 0o777,
			contents:   "test",
			errAssert:  assert.NoError,
		},
		{
			name: "success: set permission",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &filesystem.CopyFileParams{
				Src:        "/sym/src.txt",
				Dest:       "/sym/dest.txt",
				Permission: 0o644,
			},
			contents:  "test",
			errAssert: assert.NoError,
		},
		{
			name: "success: relative path current path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &filesystem.CopyFileParams{
				Src:        "./src.txt",
				Dest:       "./dest.txt",
				Permission: 0o777,
				Cwd:        "/sym",
			},
			contents:  "test",
			errAssert: assert.NoError,
		},
		{
			name: "success: relative path child path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/sym/test/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &filesystem.CopyFileParams{
				Src:        "./test/src.txt",
				Dest:       "./test/dest.txt",
				Permission: 0o777,
				Cwd:        "/sym",
			},
			contents:  "test",
			errAssert: assert.NoError,
		},
		{
			name: "success: relative path parent path",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				fs.MkdirAll("/sym/test", 0o777)
				afero.WriteFile(fs, "/sym/test2/src.txt", []byte("test"), 0o777)
				return fs
			},
			params: &filesystem.CopyFileParams{
				Src:        "../test2/src.txt",
				Dest:       "../test2/dest.txt",
				Permission: 0o777,
				Cwd:        "/sym/test",
			},
			contents:  "test",
			errAssert: assert.NoError,
		},
		{
			name: "error: source file not found",
			setup: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &filesystem.CopyFileParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			contents:  "",
			errAssert: assert.Error,
		},
		{
			name: "error: destination file not found",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0o777)
				_ = afero.WriteFile(fs, "/sym/dest.txt", []byte("test"), 0o777)
				return fs
			},
			params: &filesystem.CopyFileParams{
				Src:        "/sym/src.txt",
				Dest:       "/sym/dest.txt",
				Permission: 0o777,
			},
			contents:  "test",
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setup()
			err := filesystem.CopyFile(fs)(context.Background(), tt.params)
			tt.errAssert(t, err)

			data, _ := afero.ReadFile(fs, tt.params.Dest)
			assert.Equal(t, tt.contents, string(data))

			stat, _ := fs.Stat(tt.params.Dest)
			if stat != nil {
				assert.Equal(t, tt.params.Permission, stat.Mode())
			}
		})
	}
}
