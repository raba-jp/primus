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

func TestNewMoveFileFunction(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      filesystem.MoveFileRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: func(ctx context.Context, p *filesystem.MoveFileParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("src.txt", "dest.txt", "too many")`,
			mock: func(ctx context.Context, p *filesystem.MoveFileParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: file move failed",
			data: `test("src.txt", "dest.txt")`,
			mock: func(ctx context.Context, p *filesystem.MoveFileParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := starlark.ExecForTest("test", tt.data, filesystem.NewMoveFileFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestMoveFile(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() afero.Fs
		params    *filesystem.MoveFileParams
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
			params: &filesystem.MoveFileParams{
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
			params: &filesystem.MoveFileParams{
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
			params: &filesystem.MoveFileParams{
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
			params: &filesystem.MoveFileParams{
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
			params: &filesystem.MoveFileParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			contents:  "",
			errAssert: assert.Error,
		},
		{
			name: "error: destinatino file already exists",
			setup: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/sym/src.txt", []byte("test"), 0o777)
				_ = afero.WriteFile(fs, "/sym/dest.txt", []byte("test"), 0o777)
				return fs
			},
			params: &filesystem.MoveFileParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			contents:  "test",
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setup()
			err := filesystem.MoveFile(fs)(context.Background(), tt.params)
			tt.errAssert(t, err)

			data, _ := afero.ReadFile(fs, tt.params.Dest)
			assert.Equal(t, tt.contents, string(data))

			_, err = fs.Stat(tt.params.Src)
			assert.NotNil(t, err)
		})
	}
}
