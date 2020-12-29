package modules_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/modules"
	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewDarwinChecker(t *testing.T) {
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
					Cmd:         "uname",
					Args:        []string{"-a"},
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
							Args: exec.CmdSetStderrArgs{
								OutAnything: true,
							},
						})
						cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
							Returns: exec.CmdOutputReturns{
								Output: []byte("Darwin myMac.local 15.3.0 Darwin Kernel Version 15.3.0: Thu Dec 10 18:40:58 PST 2015; root:xnu-3248.30.4~1/RELEASE_X86_64 x86_64"),
								Err:    nil,
							},
						})
						return cmd
					},
				},
			},
			want: true,
		},
		{
			name: "fail: linux",
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "uname",
					Args:        []string{"-a"},
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
							Args: exec.CmdSetStderrArgs{
								OutAnything: true,
							},
						})
						cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
							Returns: exec.CmdOutputReturns{
								Output: []byte("Linux HostName 5.7.19-2-MANJARO #1 SMP PREEMPT Fri Aug 28 20:22:12 UTC 2020 x86_64 GNU/Linux"),
								Err:    nil,
							},
						})
						return cmd
					},
				},
			},
			want: false,
		},
		{
			name: "fail: command error",
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "uname",
					Args:        []string{"-a"},
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplySetStderrExpectation(exec.CmdSetStderrExpectation{
							Args: exec.CmdSetStderrArgs{
								OutAnything: true,
							},
						})
						cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
							Returns: exec.CmdOutputReturns{
								Output: []byte{},
								Err:    xerrors.New("dummy"),
							},
						})
						return cmd
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(exec.MockInterface)
			e.ApplyCommandContextExpectation(tt.mock)

			result := modules.NewDarwinChecker(e)(context.Background())
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestNewArchLinuxChecker(t *testing.T) {
	tests := []struct {
		name string
		mock func() afero.Fs
		want bool
	}{
		{
			name: "success",
			mock: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/etc/arch-release", []byte("Arch Linux"), 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "success: empty file",
			mock: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/etc/arch-release", []byte(""), 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "fail: not exists /etc/arch-release",
			mock: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.mock()

			result := modules.NewArchLinuxChecker(fs)(context.Background())
			assert.Equal(t, tt.want, result)
		})
	}
}
