package backend_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/internal/backend"
	"github.com/raba-jp/primus/pkg/internal/exec"
	fakeexec "github.com/raba-jp/primus/pkg/internal/exec/testing"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

func TestDarwinBackend_CheckInstall(t *testing.T) {
	tests := []struct {
		name     string
		mockExec exec.Interface
		fs       func() afero.Fs
		want     bool
	}{
		{
			name: "success: $(brew --prefix)/Celler",
			mockExec: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							Stderr: new(bytes.Buffer),
							OutputScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte("/opt/"), []byte{}, nil
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/opt/Celler/cat", []byte{}, 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "success: /opt/homebrew-cask/Caskroom",
			mockExec: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							Stderr: new(bytes.Buffer),
							OutputScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, nil
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/opt/homebrew-cask/Caskroom/cat", []byte{}, 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "success: /usr/local/Caskroom",
			mockExec: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							Stderr: new(bytes.Buffer),
							OutputScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, nil
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/usr/local/Caskroom/cat", []byte{}, 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "error: not found",
			mockExec: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							Stderr: new(bytes.Buffer),
							OutputScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, nil
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				return fs
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := backend.DarwinBackend{Exec: tt.mockExec, Fs: tt.fs()}
			if res := be.CheckInstall(context.Background(), "cat"); res != tt.want {
				t.Fatal("Fail")
			}
		})
	}
}

func TestDarwinBackend_Install(t *testing.T) {
	tests := []struct {
		name     string
		mockExec exec.Interface
		hasErr   bool
	}{
		{
			name: "success",
			mockExec: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							Stderr: new(bytes.Buffer),
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, nil
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			hasErr: false,
		},
		{
			name: "error: error occurred",
			mockExec: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							Stderr: new(bytes.Buffer),
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, xerrors.New("dummy")
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := backend.DarwinBackend{Exec: tt.mockExec, Fs: afero.NewMemMapFs()}
			if err := be.Install(context.Background(), &backend.InstallParams{Name: "base-devel", Option: "option"}); !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

func TestDarwinBackend_Uninstall(t *testing.T) {
	tests := []struct {
		name     string
		mockExec exec.Interface
		hasErr   bool
	}{
		{
			name: "success",
			mockExec: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							Stderr: new(bytes.Buffer),
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, nil
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			hasErr: false,
		},
		{
			name: "error: error occurred",
			mockExec: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							Stderr: new(bytes.Buffer),
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, xerrors.New("dummy")
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := backend.DarwinBackend{Exec: tt.mockExec, Fs: afero.NewMemMapFs()}
			if err := be.Uninstall(context.Background(), "base-devel"); !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
