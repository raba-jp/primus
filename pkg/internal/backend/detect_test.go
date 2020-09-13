package backend_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/internal/backend"
	"github.com/raba-jp/primus/pkg/internal/exec"
	fakeexec "github.com/raba-jp/primus/pkg/internal/exec/testing"
	"github.com/spf13/afero"
)

func TestDetectOS(t *testing.T) {
	tests := []struct {
		name       string
		mockStdout string
		mockFs     func() afero.Fs
		want       backend.OS
	}{
		{
			name:       "success: Darwin",
			mockStdout: "Darwin myMac.local 15.3.0 Darwin Kernel Version 15.3.0: Thu Dec 10 18:40:58 PST 2015; root:xnu-3248.30.4~1/RELEASE_X86_64 x86_64",
			mockFs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			want: backend.Darwin,
		},
		{
			name:       "success: Manjaro Linux",
			mockStdout: "Linux HostName 5.7.19-2-MANJARO #1 SMP PREEMPT Fri Aug 28 20:22:12 UTC 2020 x86_64 GNU/Linux",
			mockFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/etc/arch-release", []byte("Manjaro Linux"), 0o777)
				return fs
			},
			want: backend.Arch,
		},
		{
			name:       "fail: Unknown",
			mockStdout: "",
			mockFs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			want: backend.Unknown,
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

			if result := backend.DetectOS(execIF, tt.mockFs()); result != tt.want {
				t.Fatalf("want: %#v, got: %#v", tt.want, result)
			}
		})
	}
}

func TestDetectDarwin(t *testing.T) {
	tests := []struct {
		name       string
		mockStdout string
		want       bool
	}{
		{
			name:       "success",
			mockStdout: "Darwin myMac.local 15.3.0 Darwin Kernel Version 15.3.0: Thu Dec 10 18:40:58 PST 2015; root:xnu-3248.30.4~1/RELEASE_X86_64 x86_64",
			want:       true,
		},
		{
			name:       "fail: linux",
			mockStdout: "Linux HostName 5.7.19-2-MANJARO #1 SMP PREEMPT Fri Aug 28 20:22:12 UTC 2020 x86_64 GNU/Linux",
			want:       false,
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
			if result := backend.DetectDarwin(execIF); result != tt.want {
				t.Fatal("Fail")
			}
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

			if result := backend.DetectArchLinux(fs); result != tt.want {
				t.Fatal("Fail")
			}
		})
	}
}
