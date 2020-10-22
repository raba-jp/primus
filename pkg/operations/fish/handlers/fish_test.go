package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

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
		errAssert  assert.ErrorAssertionFunc
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
			errAssert: assert.NoError,
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
			errAssert: assert.NoError,
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
			errAssert: assert.NoError,
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
			errAssert: assert.NoError,
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
			errAssert: assert.Error,
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
			tt.errAssert(t, err)
		})
	}
}

func TestNewSetVariable__DryRun(t *testing.T) {
	tests := []struct {
		name   string
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

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestBaseBackend_FishSetPath(t *testing.T) {
	tests := []struct {
		name       string
		mockStdout string
		mockErr    error
		params     *handlers.SetPathParams
		errAssert  assert.ErrorAssertionFunc
	}{
		{
			name:       "success",
			mockStdout: "dummy",
			mockErr:    nil,
			params: &handlers.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			errAssert: assert.NoError,
		},
		{
			name:       "error",
			mockStdout: "dummy",
			mockErr:    xerrors.New("dummy"),
			params: &handlers.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			errAssert: assert.Error,
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
			tt.errAssert(t, err)
		})
	}
}

func TestNewSetPath__DryRun(t *testing.T) {
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

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
