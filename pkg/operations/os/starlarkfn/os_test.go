package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/os/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
)

func TestIsDarwin(t *testing.T) {
	tests := []struct {
		name string
		mock exec.InterfaceCommandExpectation
		want lib.Value
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
			want: lib.True,
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
			want: lib.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(exec.MockInterface)
			e.ApplyCommandExpectation(tt.mock)

			globals, err := starlark.ExecForTest("test", `v = test()`, starlarkfn.IsDarwin(e))

			assert.NoError(t, err)
			assert.Equal(t, tt.want, globals["v"])
		})
	}
}

func TestIsArchLinux(t *testing.T) {
	tests := []struct {
		name string
		mock func(fs afero.Fs)
		want lib.Value
	}{
		{
			name: "success",
			mock: func(fs afero.Fs) {
				afero.WriteFile(fs, "/etc/arch-release", []byte("Arch Linux"), 0o777)
			},
			want: lib.True,
		},
		{
			name: "success: empty file",
			mock: func(fs afero.Fs) {
				afero.WriteFile(fs, "/etc/arch-release", []byte(""), 0o777)
			},
			want: lib.True,
		},
		{
			name: "fail: not exists /etc/arch-release",
			mock: func(fs afero.Fs) {
			},
			want: lib.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			tt.mock(fs)

			globals, err := starlark.ExecForTest("test", `v = test()`, starlarkfn.IsArchLinux(fs))

			assert.NoError(t, err)
			assert.Equal(t, tt.want, globals["v"])
		})
	}
}
