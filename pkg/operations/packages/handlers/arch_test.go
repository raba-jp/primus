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
	"golang.org/x/xerrors"
)

func TestArchLinux_CheckInstall(t *testing.T) {
	tests := []struct {
		name     string
		mockExec exec.Interface
		want     bool
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
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &handlers.ArchLinux{Exec: tt.mockExec}
			res := handler.CheckInstall(context.Background(), "base-devel")
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestArchLinux_Install(t *testing.T) {
	tests := []struct {
		name      string
		mockExec  exec.Interface
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
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, xerrors.New("not installed")
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
			errAssert: assert.NoError,
		},
		{
			name: "error: install failed",
			mockExec: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							Stderr: new(bytes.Buffer),
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, xerrors.New("not installed")
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
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &handlers.ArchLinux{Exec: tt.mockExec}
			err := handler.Install(context.Background(), false, &handlers.ArchPkgInstallParams{
				Name:   "base-devel",
				Option: "option",
			})
			tt.errAssert(t, err)
		})
	}
}

func TestArchLinux_Install__dryrun(t *testing.T) {
	tests := []struct {
		name   string
		params *handlers.ArchPkgInstallParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.ArchPkgInstallParams{
				Name:   "pkg",
				Option: "option",
			},
			want: "pacman -S --noconfirm option pkg\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := &handlers.ArchLinux{}
			err := handler.Install(context.Background(), true, tt.params)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestArchLinux_Uninstall(t *testing.T) {
	tests := []struct {
		name      string
		mockExec  exec.Interface
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
							RunScript: []fakeexec.FakeAction{
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
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, xerrors.New("not installed")
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			errAssert: assert.NoError,
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
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &handlers.ArchLinux{Exec: tt.mockExec}
			err := handler.Uninstall(context.Background(), false, &handlers.ArchPkgUninstallParams{Name: "base-devel"})
			tt.errAssert(t, err)
		})
	}
}

func TestArchLinux_Uninstall__dryrun(t *testing.T) {
	tests := []struct {
		name   string
		params *handlers.ArchPkgUninstallParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.ArchPkgUninstallParams{
				Name: "pkg",
			},
			want: "pacman -R --noconfirm pkg\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := &handlers.ArchLinux{}
			err := handler.Uninstall(context.Background(), true, tt.params)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
