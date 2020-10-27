package handlers_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/darwin/handlers"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewCheckInstall(t *testing.T) {
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
			exc := new(exec.MockInterface)
			exc.ApplyCommandContextExpectation(tt.mock)

			checkInstall := handlers.NewCheckInstall(exc, tt.fs())
			res := checkInstall.Run(context.Background(), "cat")
			assert.Equal(t, tt.want, res)
		})
	}
}
