package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"golang.org/x/xerrors"
)

func TestArchLinux_CheckInstall(t *testing.T) {
	tests := []struct {
		name string
		mock exec.InterfaceCommandContextExpectation
		want bool
	}{
		{
			name: "success",
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "pacman",
					Args:        []string{"-Qg", "base-devel"},
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyRunExpectation(exec.CmdRunExpectation{
							Returns: exec.CmdRunReturns{
								Err: nil,
							},
						})
						return cmd
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectation(tt.mock)

			handler := &handlers.ArchLinux{Exec: e}
			res := handler.CheckInstall(context.Background(), "base-devel")
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestArchLinux_Install(t *testing.T) {
	tests := []struct {
		name      string
		mock      []exec.InterfaceCommandContextExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-Qg", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: xerrors.New("not installed"),
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-S", "--noconfirm", "options", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: nil,
								},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: already installed",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-Qg", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: nil,
								},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: install failed",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-Qg", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: xerrors.New("not installed"),
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-S", "--noconfirm", "options", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: xerrors.New("dummy"),
								},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectations(tt.mock)

			handler := &handlers.ArchLinux{Exec: e}
			err := handler.Install(context.Background(), false, &handlers.ArchPkgInstallParams{
				Name:   "base-devel",
				Option: "options",
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
		mock      []exec.InterfaceCommandContextExpectation
		mockExec  exec.Interface
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-Qg", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-R", "--noconfirm", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: nil,
								},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: not installed",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-Qg", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: xerrors.New("not installed"),
								},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: error occurred",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-Qg", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "pacman",
						Args:        []string{"-R", "--noconfirm", "base-devel"},
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyRunExpectation(exec.CmdRunExpectation{
								Returns: exec.CmdRunReturns{
									Err: xerrors.New("dummy"),
								},
							})
							return cmd
						},
					},
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectations(tt.mock)

			handler := &handlers.ArchLinux{Exec: e}
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
