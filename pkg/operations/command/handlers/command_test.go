package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/command/handlers"
	"golang.org/x/xerrors"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name       string
		params     *handlers.CommandParams
		mockStdout string
		mockErr    error
		errAssert  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
			},
			mockStdout: "hello world",
			mockErr:    nil,
			errAssert:  assert.NoError,
		},
		{
			name: "success: with user",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				User:    "root",
			},
			mockStdout: "hello world",
			mockErr:    nil,
			errAssert:  assert.NoError,
		},
		{
			name: "success: with cwd",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				Cwd:     "/",
			},
			mockStdout: "hello world",
			mockErr:    nil,
			errAssert:  assert.NoError,
		},
		{
			name: "success: with user and cwd",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				User:    "root",
				Cwd:     "/",
			},
			mockStdout: "hello world",
			mockErr:    nil,
			errAssert:  assert.NoError,
		},
		{
			name: "error: ",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
				User:    "root",
				Cwd:     "/",
			},
			mockStdout: "hello world",
			mockErr:    xerrors.New("dummy"),
			errAssert:  assert.Error,
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
			handler := handlers.NewCommand(execIF)

			err := handler.Command(context.Background(), false, tt.params)
			tt.errAssert(t, err)
		})
	}
}

func TestNewCommand__DryRun(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		want    string
	}{
		{
			name:    "no args",
			command: "ls",
			args:    []string{},
			want:    "ls \n",
		},
		{
			name:    "add option",
			command: "ls",
			args:    []string{"-al"},
			want:    "ls -al \n",
		},
		{
			name:    "with double quote",
			command: "ls",
			args:    []string{"-al", "\"go.mod\""},
			want:    "ls -al \"go.mod\" \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			handler := handlers.NewCommand(nil)
			err := handler.Command(context.Background(), true, &handlers.CommandParams{
				CmdName: tt.command,
				CmdArgs: tt.args,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
