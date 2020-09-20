package functions_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/exec"
	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
)

func TestIsDarwin(t *testing.T) {
	tests := []struct {
		name       string
		mockStdout string
		want       starlark.Value
	}{
		{
			name:       "success",
			mockStdout: "Darwin myMac.local 15.3.0 Darwin Kernel Version 15.3.0: Thu Dec 10 18:40:58 PST 2015; root:xnu-3248.30.4~1/RELEASE_X86_64 x86_64",
			want:       starlark.True,
		},
		{
			name:       "fail: linux",
			mockStdout: "Linux HostName 5.7.19-2-MANJARO #1 SMP PREEMPT Fri Aug 28 20:22:12 UTC 2020 x86_64 GNU/Linux",
			want:       starlark.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execIF := &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							OutputScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte(tt.mockStdout), []byte{}, nil
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			}
			predeclared := starlark.StringDict{
				"is_darwin": starlark.NewBuiltin("is_darwin", functions.IsDarwin(execIF)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			globals, err := starlark.ExecFile(thread, "test.star", `v = is_darwin()`, predeclared)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if globals["v"] != tt.want {
				t.Errorf("want: %v, got: %v", tt.want, globals["v"])
			}
		})
	}
}

func TestIsArchLinux(t *testing.T) {
	tests := []struct {
		name string
		mock func(fs afero.Fs)
		want starlark.Value
	}{
		{
			name: "success",
			mock: func(fs afero.Fs) {
				afero.WriteFile(fs, "/etc/arch-release", []byte("Arch Linux"), 0o777)
			},
			want: starlark.True,
		},
		{
			name: "success: empty file",
			mock: func(fs afero.Fs) {
				afero.WriteFile(fs, "/etc/arch-release", []byte(""), 0o777)
			},
			want: starlark.True,
		},
		{
			name: "fail: not exists /etc/arch-release",
			mock: func(fs afero.Fs) {
			},
			want: starlark.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			tt.mock(fs)

			predeclared := starlark.StringDict{
				"is_arch_linux": starlark.NewBuiltin("is_arch_linux", functions.IsArchLinux(fs)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			globals, err := starlark.ExecFile(thread, "test.star", `v = is_arch_linux()`, predeclared)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if globals["v"] != tt.want {
				t.Errorf("want: %v, got: %v", tt.want, globals["v"])
			}
		})
	}
}
