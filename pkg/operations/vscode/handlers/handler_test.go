package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"
	"github.com/raba-jp/primus/pkg/operations/vscode/handlers"
	"golang.org/x/xerrors"
)

func TestNewInstallExtension(t *testing.T) {
	tests := []struct {
		name       string
		mockStdout string
		mockErr    error
		params     *handlers.InstallExtensionParams
		hasErr     bool
	}{
		{
			name:       "success: without version",
			mockStdout: "dummy",
			mockErr:    nil,
			params: &handlers.InstallExtensionParams{
				Name:    "test",
				Version: "",
			},
			hasErr: false,
		},
		{
			name:       "success: with version",
			mockStdout: "dummy",
			mockErr:    nil,
			params: &handlers.InstallExtensionParams{
				Name:    "test",
				Version: "1.0",
			},
			hasErr: false,
		},
		{
			name:       "error",
			mockStdout: "dummy",
			mockErr:    xerrors.New("dummy"),
			params: &handlers.InstallExtensionParams{
				Name:    "test",
				Version: "",
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execIF := &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte(tt.mockStdout), []byte{}, tt.mockErr
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			}

			handler := handlers.NewInstallExtension(execIF)
			err := handler.InstallExtension(context.Background(), false, tt.params)
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

func TestNewInstallExtension__DryRun(t *testing.T) {
	tests := []struct {
		name   string
		params *handlers.InstallExtensionParams
		want   string
	}{
		{
			name: "success: without version",
			params: &handlers.InstallExtensionParams{
				Name:    "test",
				Version: "",
			},
			want: "code --install-extension test\n",
		},
		{
			name: "success: with version",
			params: &handlers.InstallExtensionParams{
				Name:    "test",
				Version: "1.0",
			},
			want: "code --install-extension test@1.0\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.NewInstallExtension(nil)
			err := handler.InstallExtension(context.Background(), true, tt.params)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
