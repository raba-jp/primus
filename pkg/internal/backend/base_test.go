package backend_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/internal/backend"
	"github.com/raba-jp/primus/pkg/internal/exec"
	fakeexec "github.com/raba-jp/primus/pkg/internal/exec/testing"
	"github.com/raba-jp/primus/pkg/internal/handlers"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

type MockRoundTripper struct {
	http.RoundTripper
	Fn func(req *http.Request) *http.Response
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Fn(req), nil
}

func MockHttpClient(fn func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{Fn: fn},
	}
}

func TestBaseBackend_Command(t *testing.T) {
	tests := []struct {
		name       string
		params     *handlers.CommandParams
		mockStdout string
		mockErr    error
		hasErr     bool
	}{
		{
			name: "success",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
			},
			mockStdout: "hello world",
			mockErr:    nil,
			hasErr:     false,
		},
		{
			name: "success: with user",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				User:    "root",
			},
			mockStdout: "hello world",
			mockErr:    nil,
			hasErr:     false,
		},
		{
			name: "success: with cwd",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				Cwd:     "/",
			},
			mockStdout: "hello world",
			mockErr:    nil,
			hasErr:     false,
		},
		{
			name: "success: with user and cwd",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				User:    "root",
				Cwd:     "/",
			},
			mockStdout: "hello world",
			mockErr:    nil,
			hasErr:     false,
		},
		{
			name: "error: ",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				User:    "root",
				Cwd:     "/",
			},
			mockStdout: "hello world",
			mockErr:    xerrors.New("dummy"),
			hasErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execIF := &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte(tt.mockStdout), []byte{}, tt.mockErr
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			}
			be := backend.BaseBackend{Exec: execIF}

			err := be.Command(context.Background(), false, tt.params)
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

func TestBaseBackend_Command__DryRun(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		want    string
	}{
		{
			name:    "no args",
			command: "ls",
			args:    []string{},
			want:    "ls \n",
		},
		{
			name:    "add option",
			command: "ls",
			args:    []string{"-al"},
			want:    "ls -al \n",
		},
		{
			name:    "with double quote",
			command: "ls",
			args:    []string{"-al", "\"go.mod\""},
			want:    "ls -al \"go.mod\" \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			be := backend.BaseBackend{}
			err := be.Command(context.Background(), true, &handlers.CommandParams{
				CmdName: tt.command,
				CmdArgs: tt.args,
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBaseBackend_FileCopy(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() afero.Fs
		params     *handlers.FileCopyParams
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
			params: &handlers.FileCopyParams{
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
			params: &handlers.FileCopyParams{
				Src:        "/sym/src.txt",
				Dest:       "/sym/dest.txt",
				Permission: 0o644,
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "error: source file not found",
			setup: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &handlers.FileCopyParams{
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
			be := &backend.BaseBackend{Fs: fs}
			err := be.FileCopy(context.Background(), false, tt.params)
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

			be := &backend.BaseBackend{}
			err := be.FileCopy(context.Background(), true, &handlers.FileCopyParams{
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

func TestBaseBackend_FileMove(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() afero.Fs
		params   *handlers.FileMoveParams
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
			params: &handlers.FileMoveParams{
				Src:  "/sym/src.txt",
				Dest: "/sym/dest.txt",
			},
			contents: "test",
			hasErr:   false,
		},
		{
			name: "error: source file not found",
			setup: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &handlers.FileMoveParams{
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

			be := backend.BaseBackend{Fs: fs}
			err := be.FileMove(context.Background(), false, tt.params)
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

func TestBaseBackend_FileMove__DryRun(t *testing.T) {
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

			be := &backend.BaseBackend{}
			err := be.FileMove(context.Background(), true, &handlers.FileMoveParams{
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

func TestBackend_HTTPRequest(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		path     string
		contents string
		httpMock func(req *http.Request) *http.Response
	}{
		{
			name:     "success",
			url:      "https://example.com/",
			path:     "/sym/test.txt",
			contents: "test file",
			httpMock: func(req *http.Request) *http.Response {
				buf := bytes.NewBufferString("test file")
				body := ioutil.NopCloser(buf)
				return &http.Response{
					Body: body,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			be := backend.BaseBackend{Fs: fs, Client: MockHttpClient(tt.httpMock)}
			err := be.HTTPRequest(context.Background(), false, &handlers.HTTPRequestParams{
				URL:  tt.url,
				Path: tt.path,
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
			data, err := afero.ReadFile(fs, tt.path)
			if err != nil {
				t.Fatalf("file read failed: %s: %v", tt.path, err)
			}
			if diff := cmp.Diff(tt.contents, string(data)); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBaseBackend_HTTPRequest__DryRun(t *testing.T) {
	tests := []struct {
		name string
		url  string
		path string
		want string
	}{
		{
			name: "success",
			url:  "https://example.com",
			path: "/sym/output.txt",
			want: "curl -Lo /sym/output.txt https://example.com\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			be := backend.BaseBackend{}
			err := be.HTTPRequest(context.Background(), true, &handlers.HTTPRequestParams{
				URL:  tt.url,
				Path: tt.path,
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

func TestBackend_Symlink(t *testing.T) {
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

			be := backend.BaseBackend{Fs: fs}
			err := be.Symlink(context.Background(), false, &handlers.SymlinkParams{
				Src:  tt.src,
				Dest: tt.dest,
			})
			if err != nil {
				t.Fatalf("%v", err)
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

			be := backend.BaseBackend{Fs: fs}
			_ = be.Symlink(context.Background(), false, &handlers.SymlinkParams{
				Src:  tt.src,
				Dest: tt.dest,
			})

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

			be := backend.BaseBackend{Fs: fs}
			_ = be.Symlink(context.Background(), false, &handlers.SymlinkParams{
				Src:  tt.src,
				Dest: tt.dest,
			})
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

func TestBaseBackend_Symlink__DryRun(t *testing.T) {
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
			want: "ln -s /sym/src.txt /sym/dest.txt\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			be := &backend.BaseBackend{}
			err := be.Symlink(context.Background(), true, &handlers.SymlinkParams{
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
