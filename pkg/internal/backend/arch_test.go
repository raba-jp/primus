package backend_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/internal/backend"
	"github.com/raba-jp/primus/pkg/internal/exec"
	fakeexec "github.com/raba-jp/primus/pkg/internal/exec/testing"
	"golang.org/x/xerrors"
)

func TestArchLinuxBackend_CheckInstall(t *testing.T) {
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
			be := backend.ArchLinuxBackend{Exec: tt.mockExec}
			if res := be.CheckInstall(context.Background(), "base-devel"); res != tt.want {
				t.Fatal("Fail")
			}
		})
	}
}

func TestArchLinuxBackend_Install(t *testing.T) {
	tests := []struct {
		name     string
		mockExec exec.Interface
		hasErr   bool
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
			hasErr: false,
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
									return []byte{}, []byte{}, xerrors.New("dummy")
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := backend.ArchLinuxBackend{Exec: tt.mockExec}
			if err := be.Install(context.Background(), false, &backend.InstallParams{Name: "base-devel", Option: "option"}); !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

func TestArchLinuxBackend_Uninstall(t *testing.T) {
	tests := []struct {
		name     string
		mockExec exec.Interface
		hasErr   bool
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
			hasErr: false,
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
									return []byte{}, []byte{}, xerrors.New("dummy")
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := backend.ArchLinuxBackend{Exec: tt.mockExec}
			if err := be.Uninstall(context.Background(), false, "base-devel"); !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
