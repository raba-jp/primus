package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"
	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

func TestNewDarwinPkgCheckInstall(t *testing.T) {
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
			handler := &handlers.Darwin{Exec: tt.mockExec, Fs: tt.fs()}
			res := handler.CheckInstall(context.Background(), "cat")
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestNewDarwinPkgInstall(t *testing.T) {
	tests := []struct {
		name      string
		mockExec  exec.Interface
		fs        func() afero.Fs
		params    *handlers.DarwinPkgInstallParams
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
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
			params: &handlers.DarwinPkgInstallParams{
				Name:   "pkg",
				Option: "options",
			},
			fs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: already installed",
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
			params: &handlers.DarwinPkgInstallParams{
				Name:   "pkg",
				Option: "options",
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/usr/local/Caskroom/pkg", []byte{}, 0o777)
				return fs
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: install package failed",
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
			params: &handlers.DarwinPkgInstallParams{
				Name:   "pkg",
				Option: "options",
			},
			fs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.Darwin{Exec: tt.mockExec, Fs: tt.fs()}
			err := handler.Install(context.Background(), false, tt.params)
			tt.errAssert(t, err)
		})
	}
}

func TestTestDarwinPkgInstall__dryrun(t *testing.T) {
	tests := []struct {
		name   string
		params *handlers.DarwinPkgInstallParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.DarwinPkgInstallParams{
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

			handler := handlers.Darwin{}
			err := handler.Install(context.Background(), true, tt.params)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestDarwin_Uninstall(t *testing.T) {
	tests := []struct {
		name      string
		mockExec  exec.Interface
		fs        func() afero.Fs
		params    *handlers.DarwinPkgUninstallParams
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
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
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/usr/local/Caskroom/pkg", []byte{}, 0o777)
				return fs
			},
			params:    &handlers.DarwinPkgUninstallParams{Name: "pkg"},
			errAssert: assert.NoError,
		},
		{
			name: "success: not installed",
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
			fs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params:    &handlers.DarwinPkgUninstallParams{Name: "pkg"},
			errAssert: assert.NoError,
		},
		{
			name: "error: uninstall failed",
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
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/usr/local/Caskroom/pkg", []byte{}, 0o777)
				return fs
			},
			params:    &handlers.DarwinPkgUninstallParams{Name: "pkg"},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.Darwin{Exec: tt.mockExec, Fs: tt.fs()}
			err := handler.Uninstall(context.Background(), false, tt.params)
			tt.errAssert(t, err)
		})
	}
}

func TestDarwin_Uninstall__dryrun(t *testing.T) {
	tests := []struct {
		name   string
		params *handlers.DarwinPkgUninstallParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.DarwinPkgUninstallParams{
				Name: "pkg",
			},
			want: "brew uninstall pkg\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.Darwin{}
			err := handler.Uninstall(context.Background(), true, tt.params)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
