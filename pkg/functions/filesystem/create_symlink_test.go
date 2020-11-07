package filesystem_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/raba-jp/primus/pkg/functions/filesystem"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewCreateSymlinkFunction(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      filesystem.CreateSymlinkRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: func(ctx context.Context, p *filesystem.CreateSymlinkParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("/sym/src.txt", "/sym/dest.txt", "too many")`,
			mock: func(ctx context.Context, p *filesystem.CreateSymlinkParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: create symlink failed ",
			data: `test("/sym/src.txt", "/sym/dest.txt")`,
			mock: func(ctx context.Context, p *filesystem.CreateSymlinkParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := starlark.ExecForTest("test", tt.data, filesystem.NewCreateSymlinkFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestCreateSymlink(t *testing.T) {
	wd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(wd)
	}()
	fs := afero.NewOsFs()

	workDir, err := afero.TempDir(fs, "", "primus_symlink_test")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = fs.RemoveAll(workDir)
	}()

	tests := []struct {
		name string
		src  string
		dest string
	}{
		{
			name: "success",
			src:  filepath.Join(workDir, "src.txt"),
			dest: filepath.Join(workDir, "dest.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = afero.WriteFile(fs, tt.src, []byte("test file"), 0777)

			err := filesystem.CreateSymlink(fs)(context.Background(), &filesystem.CreateSymlinkParams{
				Src:  tt.src,
				Dest: tt.dest,
			})
			assert.NoError(t, err)

			lst, _ := fs.(afero.Lstater)
			_, ok, err := lst.LstatIfPossible(tt.dest)
			if !ok {
				if err != nil {
					t.Fatalf(": %v", err)
				} else {
					t.Fatal("Error calling lstat(not link)")
				}
			}
		})
	}
}

func TestSymlink_AlreadyExistsFile(t *testing.T) {
	wd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(wd)
	}()
	fs := afero.NewOsFs()

	workDir, err := afero.TempDir(fs, "", "symlink_iac")
	assert.NoError(t, err)

	defer func() {
		_ = fs.RemoveAll(workDir)
	}()

	tests := []struct {
		src  string
		dest string
	}{
		{
			src:  filepath.Join(workDir, "src.txt"),
			dest: filepath.Join(workDir, "dest.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s to %s", tt.src, tt.dest), func(t *testing.T) {

			_ = afero.WriteFile(fs, tt.src, []byte("test file"), 0777)
			_ = afero.WriteFile(fs, tt.dest, []byte("test file"), 0777)

			_ = filesystem.CreateSymlink(fs)(context.Background(), &filesystem.CreateSymlinkParams{
				Src:  tt.src,
				Dest: tt.dest,
			})

			_, err = fs.Stat(tt.dest)
			assert.NoError(t, err)
		})
	}
}

func TestSymlink_AlreadyExistsSymlink(t *testing.T) {
	wd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(wd)
	}()
	fs := afero.NewOsFs()

	workDir, err := afero.TempDir(fs, "", "symlink_iac")
	assert.NoError(t, err)

	defer func() {
		_ = fs.RemoveAll(workDir)
	}()

	tests := []struct {
		src  string
		dest string
	}{
		{
			src:  filepath.Join(workDir, "src.txt"),
			dest: filepath.Join(workDir, "dest.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s to %s", tt.src, tt.dest), func(t *testing.T) {
			_ = afero.WriteFile(fs, tt.src, []byte("test file"), 0777)
			another := filepath.Join(workDir, "another.txt")
			_ = afero.WriteFile(fs, another, []byte("another test file"), 0777)
			l, _ := fs.(afero.Linker)
			_ = l.SymlinkIfPossible(another, tt.dest)

			_ = filesystem.CreateSymlink(fs)(context.Background(), &filesystem.CreateSymlinkParams{
				Src:  tt.src,
				Dest: tt.dest,
			})
			data, err := afero.ReadFile(fs, tt.dest)
			assert.NoError(t, err)
			assert.Equal(t, "another test file", string(data))
		})
	}
}
