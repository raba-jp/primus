package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"
	"github.com/raba-jp/primus/pkg/operations/fish/handlers"
	"golang.org/x/xerrors"
)

func TestNewSetVariable(t *testing.T) {
	tests := []struct {
		name       string
		mockStdout string
		mockErr    error
		params     *handlers.SetVariableParams
		hasErr     bool
	}{
		{
			name:       "success: scope universal",
			mockStdout: "dummy",
			mockErr:    nil,
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.UniversalScope,
				Export: true,
			},
			hasErr: false,
		},
		{
			name:       "success: scope global",
			mockStdout: "dummy",
			mockErr:    nil,
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.GlobalScope,
				Export: true,
			},
			hasErr: false,
		},
		{
			name:       "success: scope local",
			mockStdout: "dummy",
			mockErr:    nil,
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.LocalScope,
				Export: true,
			},
			hasErr: false,
		},
		{
			name:       "success: no export",
			mockStdout: "dummy",
			mockErr:    nil,
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.LocalScope,
				Export: false,
			},
			hasErr: false,
		},
		{
			name:       "error",
			mockStdout: "dummy",
			mockErr:    xerrors.New("dummy"),
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.UniversalScope,
				Export: true,
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

			handler := handlers.NewSetVariable(execIF)
			err := handler.SetVariable(context.Background(), false, tt.params)
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

func TestBaseBackend_FishSetVariable__DryRun(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		params *handlers.SetVariableParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  handlers.UniversalScope,
				Export: true,
			},
			want: "fish --command 'set --universal --export GOPATH $HOME/go'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.NewSetVariable(nil)
			err := handler.SetVariable(context.Background(), true, tt.params)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBaseBackend_FishSetPath(t *testing.T) {
	tests := []struct {
		name       string
		mockStdout string
		mockErr    error
		params     *handlers.SetPathParams
		hasErr     bool
	}{
		{
			name:       "success",
			mockStdout: "dummy",
			mockErr:    nil,
			params: &handlers.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			hasErr: false,
		},
		{
			name:       "error",
			mockStdout: "dummy",
			mockErr:    xerrors.New("dummy"),
			params: &handlers.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
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

			handler := handlers.NewSetPath(execIF)
			err := handler.SetPath(context.Background(), false, tt.params)
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

func TestBaseBackend_FishSetPath__DryRun(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		params *handlers.SetPathParams
		want   string
	}{
		{
			name: "success",
			params: &handlers.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			want: "fish --command 'set --universal fish_user_paths $GOPATH/bin $HOME/.bin'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.NewSetPath(nil)
			err := handler.SetPath(context.Background(), true, tt.params)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
