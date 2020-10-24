package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

func TestNewDarwinPkgCheckInstall(t *testing.T) {
	tests := []struct {
		name string
		mock exec.InterfaceCommandContextExpectation
		fs   func() afero.Fs
		want bool
	}{
		{
			name: "success: $(brew --prefix)/Celler",
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything:  true,
					CmdAnything:  true,
					ArgsAnything: true,
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
							Returns: exec.CmdOutputReturns{
								Output: []byte("/opt"),
								Err:    nil,
							},
						})
						return cmd
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
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything:  true,
					CmdAnything:  true,
					ArgsAnything: true,
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
							Returns: exec.CmdOutputReturns{
								Output: []byte("/opt"),
								Err:    nil,
							},
						})
						return cmd
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
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything:  true,
					CmdAnything:  true,
					ArgsAnything: true,
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
							Returns: exec.CmdOutputReturns{
								Output: []byte("/opt"),
								Err:    nil,
							},
						})
						return cmd
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
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything:  true,
					CmdAnything:  true,
					ArgsAnything: true,
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
							Returns: exec.CmdOutputReturns{
								Output: []byte("/opt"),
								Err:    nil,
							},
						})
						return cmd
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
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectation(tt.mock)

			handler := &handlers.Darwin{Exec: e, Fs: tt.fs()}
			res := handler.CheckInstall(context.Background(), "cat")
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestNewDarwinPkgInstall(t *testing.T) {
	tests := []struct {
		name      string
		mock      []exec.InterfaceCommandContextExpectation
		fs        func() afero.Fs
		params    *handlers.DarwinPkgInstallParams
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything:  true,
						CmdAnything:  true,
						ArgsAnything: true,
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte("/opt"),
									Err:    nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "brew",
						Args:        []string{"install", "options", "pkg"},
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
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything:  true,
						CmdAnything:  true,
						ArgsAnything: true,
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte("/opt"),
									Err:    nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything:  true,
						CmdAnything:  true,
						ArgsAnything: true,
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
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything:  true,
						CmdAnything:  true,
						ArgsAnything: true,
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte("/opt"),
									Err:    nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "brew",
						Args:        []string{"install", "options", "pkg"},
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
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectations(tt.mock)

			handler := handlers.Darwin{Exec: e, Fs: tt.fs()}
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
		mock      []exec.InterfaceCommandContextExpectation
		fs        func() afero.Fs
		params    *handlers.DarwinPkgUninstallParams
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything:  true,
						CmdAnything:  true,
						ArgsAnything: true,
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte("/opt"),
									Err:    nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "brew",
						Args:        []string{"uninstall", "pkg"},
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
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything:  true,
						CmdAnything:  true,
						ArgsAnything: true,
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte("/opt"),
									Err:    nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "brew",
						Args:        []string{"uninstall", "pkg"},
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
			fs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params:    &handlers.DarwinPkgUninstallParams{Name: "pkg"},
			errAssert: assert.NoError,
		},
		{
			name: "error: uninstall failed",
			mock: []exec.InterfaceCommandContextExpectation{
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything:  true,
						CmdAnything:  true,
						ArgsAnything: true,
					},
					Returns: exec.InterfaceCommandContextReturns{
						Cmd: func() exec.Cmd {
							cmd := new(exec.MockCmd)
							cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
								Returns: exec.CmdOutputReturns{
									Output: []byte("/opt"),
									Err:    nil,
								},
							})
							return cmd
						},
					},
				},
				{
					Args: exec.InterfaceCommandContextArgs{
						CtxAnything: true,
						Cmd:         "brew",
						Args:        []string{"uninstall", "pkg"},
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
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectations(tt.mock)
			handler := handlers.Darwin{Exec: e, Fs: tt.fs()}
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
