package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"
	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

func TestDarwin_CheckInstall(t *testing.T) {
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
			handler := handlers.Darwin{Exec: tt.mockExec, Fs: tt.fs()}
			if res := handler.CheckInstall(context.Background(), "cat"); res != tt.want {
				t.Fatal("Fail")
			}
		})
	}
}

func TestDarwin_Install(t *testing.T) {
	tests := []struct {
		name     string
		mockExec exec.Interface
		params   *handlers.InstallParams
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
			params: &handlers.InstallParams{
				Name:   "pkg",
				Option: "options",
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
			params: &handlers.InstallParams{
				Name:   "pkg",
				Option: "options",
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.Darwin{Exec: tt.mockExec}
			if err := handler.Install(context.Background(), false, tt.params); !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

func TestDarwin_Install__dryrun(t *testing.T) {
	tests := []struct {
		name   string
		params *handlers.InstallParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.InstallParams{
				Name:   "pkg",
				Option: "option",
			},
			want: "brew install option pkg\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := &handlers.Darwin{}
			err := handler.Install(context.Background(), true, tt.params)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestDarwin_Uninstall(t *testing.T) {
	tests := []struct {
		name     string
		mockExec exec.Interface
		params   *handlers.UninstallParams
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
			params: &handlers.UninstallParams{Name: "pkg"},
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
			params: &handlers.UninstallParams{Name: "pkg"},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.Darwin{Exec: tt.mockExec}
			if err := handler.Uninstall(context.Background(), false, tt.params); !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

func TestDarwin_Uninstall__dryrun(t *testing.T) {
	tests := []struct {
		name   string
		params *handlers.UninstallParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.UninstallParams{
				Name: "pkg",
			},
			want: "brew uninstall pkg\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := &handlers.Darwin{}
			err := handler.Uninstall(context.Background(), true, tt.params)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
