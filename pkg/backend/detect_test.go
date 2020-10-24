package backend_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDetectOS(t *testing.T) {
	tests := []struct {
		name       string
		mock       exec.InterfaceCommandExpectation
		mockStdout string
		mockFs     func() afero.Fs
		want       backend.OS
	}{
		{
			name: "success: Darwin",
			mock: exec.InterfaceCommandExpectation{
				Args: exec.InterfaceCommandArgs{
					CmdAnything:  true,
					ArgsAnything: true,
				},
				Returns: exec.InterfaceCommandReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
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
			mockFs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			want: backend.Darwin,
		},
		{
			name: "success: Manjaro Linux",
			mock: exec.InterfaceCommandExpectation{
				Args: exec.InterfaceCommandArgs{
					CmdAnything:  true,
					ArgsAnything: true,
				},
				Returns: exec.InterfaceCommandReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
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
			mockFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/etc/arch-release", []byte("Manjaro Linux"), 0o777)
				return fs
			},
			want: backend.Arch,
		},
		{
			name: "fail: Unknown",
			mock: exec.InterfaceCommandExpectation{
				Args: exec.InterfaceCommandArgs{
					CmdAnything:  true,
					ArgsAnything: true,
				},
				Returns: exec.InterfaceCommandReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyOutputExpectation(exec.CmdOutputExpectation{
							Returns: exec.CmdOutputReturns{
								Output: []byte{},
								Err:    nil,
							},
						})
						return cmd
					},
				},
			},
			mockFs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			want: backend.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(exec.MockInterface)
			e.ApplyCommandExpectation(tt.mock)

			result := backend.DetectOS(e, tt.mockFs())
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestDetectDarwin(t *testing.T) {
	tests := []struct {
		name string
		mock exec.InterfaceCommandExpectation
		want bool
	}{
		{
			name: "success",
			mock: exec.InterfaceCommandExpectation{
				Args: exec.InterfaceCommandArgs{
					CmdAnything:  true,
					ArgsAnything: true,
				},
				Returns: exec.InterfaceCommandReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
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
			mock: exec.InterfaceCommandExpectation{
				Args: exec.InterfaceCommandArgs{
					CmdAnything:  true,
					ArgsAnything: true,
				},
				Returns: exec.InterfaceCommandReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(exec.MockInterface)
			e.ApplyCommandExpectation(tt.mock)

			result := backend.DetectDarwin(e)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestDetectManjaroLinux(t *testing.T) {
	tests := []struct {
		name string
		mock func(fs afero.Fs)
		want bool
	}{
		{
			name: "success",
			mock: func(fs afero.Fs) {
				afero.WriteFile(fs, "/etc/arch-release", []byte("Arch Linux"), 0o777)
			},
			want: true,
		},
		{
			name: "success: empty file",
			mock: func(fs afero.Fs) {
				afero.WriteFile(fs, "/etc/arch-release", []byte(""), 0o777)
			},
			want: true,
		},
		{
			name: "fail: not exists /etc/arch-release",
			mock: func(fs afero.Fs) {
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			tt.mock(fs)

			result := backend.DetectArchLinux(fs)
			assert.Equal(t, tt.want, result)
		})
	}
}
