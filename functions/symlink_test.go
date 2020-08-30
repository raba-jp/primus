package functions_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/raba-jp/primus/functions"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
)

func TestSymlink(t *testing.T) {
	wd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(wd)
	}()
	fs := afero.NewOsFs()

	workDir, err := afero.TempDir(fs, "", "symlink_iac")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = fs.RemoveAll(workDir)
	}()

	tests := []struct {
		data string
		src  string
		dest string
	}{
		{
			data: fmt.Sprintf(
				`symlink(src="%s", dest="%s")`,
				filepath.Join(workDir, "src.txt"),
				filepath.Join(workDir, "dest.txt"),
			),
			src:  filepath.Join(workDir, "src.txt"),
			dest: filepath.Join(workDir, "dest.txt"),
		},
		{
			data: fmt.Sprintf(
				`symlink("%s", "%s")`,
				filepath.Join(workDir, "src1.txt"),
				filepath.Join(workDir, "dest1.txt"),
			),
			src:  filepath.Join(workDir, "src1.txt"),
			dest: filepath.Join(workDir, "dest1.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			_ = afero.WriteFile(fs, tt.src, []byte("test file"), 0777)

			predeclared := starlark.StringDict{
				"symlink": starlark.NewBuiltin("symlink", functions.Symlink(context.Background(), fs)),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}

			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if err != nil {
				t.Fatal(err)
			}

			lst, _ := fs.(afero.Lstater)
			_, ok, err := lst.LstatIfPossible(tt.dest)
			if !ok {
				if err != nil {
					t.Fatalf("Error calling lstat: %v", err)
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
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = fs.RemoveAll(workDir)
	}()

	tests := []struct {
		data string
		src  string
		dest string
	}{
		{
			data: fmt.Sprintf(
				`symlink(src="%s", dest="%s")`,
				filepath.Join(workDir, "src.txt"),
				filepath.Join(workDir, "dest.txt"),
			),
			src:  filepath.Join(workDir, "src.txt"),
			dest: filepath.Join(workDir, "dest.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			_ = afero.WriteFile(fs, tt.src, []byte("test file"), 0777)
			_ = afero.WriteFile(fs, tt.dest, []byte("test file"), 0777)

			predeclared := starlark.StringDict{
				"symlink": starlark.NewBuiltin("symlink", functions.Symlink(context.Background(), fs)),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}

			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if err != nil {
				t.Fatal(err)
			}

			_, err = fs.Stat(tt.dest)
			if err != nil {
				t.Fatalf("%v", err)
			}
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
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = fs.RemoveAll(workDir)
	}()

	tests := []struct {
		data string
		src  string
		dest string
	}{
		{
			data: fmt.Sprintf(
				`symlink(src="%s", dest="%s")`,
				filepath.Join(workDir, "src.txt"),
				filepath.Join(workDir, "dest.txt"),
			),
			src:  filepath.Join(workDir, "src.txt"),
			dest: filepath.Join(workDir, "dest.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			_ = afero.WriteFile(fs, tt.src, []byte("test file"), 0777)
			another := filepath.Join(workDir, "another.txt")
			_ = afero.WriteFile(fs, another, []byte("another test file"), 0777)
			l, _ := fs.(afero.Linker)
			_ = l.SymlinkIfPossible(another, tt.dest)

			predeclared := starlark.StringDict{
				"symlink": starlark.NewBuiltin("symlink", functions.Symlink(context.Background(), fs)),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}

			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if err != nil {
				t.Fatal(err)
			}

			data, err := afero.ReadFile(fs, tt.dest)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if string(data) != "another test file" {
				t.Fatal("unexpected symlink")
			}
		})
	}
}
