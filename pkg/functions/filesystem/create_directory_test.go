package filesystem_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/functions/filesystem"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewCreateDirectoryFunction(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      filesystem.CreateDirectoryRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(path="/sym/test", permission=0o777)`,
			mock: func(ctx context.Context, p *filesystem.CreateDirectoryParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: relative path",
			data: `test(path="test", permission=0o777)`,
			mock: func(ctx context.Context, p *filesystem.CreateDirectoryParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: without permission",
			data: `test(path="/sym/test")`,
			mock: func(ctx context.Context, p *filesystem.CreateDirectoryParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with cwd",
			data: `test(path="test", cwd="/sym")`,
			mock: func(ctx context.Context, p *filesystem.CreateDirectoryParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("/sym/test", 0o644, "", "too many")`,
			mock: func(ctx context.Context, p *filesystem.CreateDirectoryParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to create directory",
			data: `test("/sym/test")`,
			mock: func(ctx context.Context, p *filesystem.CreateDirectoryParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := starlark.ExecForTest(
				"test",
				tt.data,
				filesystem.NewCreateDirectoryFunction(tt.mock),
			)
			tt.errAssert(t, err)
		})
	}
}

func TestCreateDirectory(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() afero.Fs
		params    *filesystem.CreateDirectoryParams
		got       string
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			setup: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &filesystem.CreateDirectoryParams{
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
			params: &filesystem.CreateDirectoryParams{
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
			params: &filesystem.CreateDirectoryParams{
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
			params: &filesystem.CreateDirectoryParams{
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
			params: &filesystem.CreateDirectoryParams{
				Path:       "../test2/test3",
				Permission: 0o644,
				Cwd:        "/sym/test",
			},
			got:       "/sym/test2/test3",
			errAssert: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setup()
			err := filesystem.CreateDirectory(fs)(context.Background(), tt.params)
			tt.errAssert(t, err)

			stat, err := fs.Stat(tt.got)
			assert.NoErrorf(t, err, "create directory failed")

			isdir := !stat.IsDir()
			assert.Falsef(t, isdir, tt.got+"is not directory")
		})
	}
}
