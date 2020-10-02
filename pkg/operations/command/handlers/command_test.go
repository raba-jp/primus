package handlers_test

import (
	"bytes"
	"context"
	"testing"

	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"

	"github.com/google/go-cmp/cmp"
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
		hasErr     bool
	}{
		{
			name: "success",
			params: &handlers.CommandParams{
				CmdName: "echo",
				CmdArgs: []string{"hello", "world"},
			},
			mockStdout: "hello world",
			mockErr:    nil,
			hasErr:     false,
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
			hasErr:     false,
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
			hasErr:     false,
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
			hasErr:     false,
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
			hasErr:     true,
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
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
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
			if err != nil {
				t.Fatalf("%v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
