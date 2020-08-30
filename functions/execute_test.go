package functions_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/functions"
	"go.starlark.net/starlark"
	"k8s.io/utils/exec"
	fakeexec "k8s.io/utils/exec/testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		data     string
		wantCmd  string
		wantArgs []string
		hasErr   bool
	}{
		{
			data:     `execute(name="echo", args=["hello", "world"])`,
			wantCmd:  "echo",
			wantArgs: []string{"hello", "world"},
			hasErr:   false,
		},
		{
			data:     `execute(name="echo", args=[1])`,
			wantCmd:  "echo",
			wantArgs: []string{"1"},
			hasErr:   false,
		},
		{
			data:     `execute(name="echo", args=[False, True])`,
			wantCmd:  "echo",
			wantArgs: []string{"false", "true"},
			hasErr:   false,
		},
		{
			data:     `execute(name="echo", args=[1.111])`,
			wantCmd:  "",
			wantArgs: nil,
			hasErr:   true,
		},
		{
			data:     `execute("echo", ["hello", "world"])`,
			wantCmd:  "echo",
			wantArgs: []string{"hello", "world"},
			hasErr:   false,
		},
		{
			data:     `execute("echo")`,
			wantCmd:  "echo",
			wantArgs: []string{},
			hasErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			var gotCmd string
			var gotArgs []string
			fexec := fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						gotCmd = cmd
						gotArgs = args

						fake := &fakeexec.FakeCmd{
							Stdout: new(bytes.Buffer),
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte{}, []byte{}, nil
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			}
			predeclared := starlark.StringDict{
				"execute": starlark.NewBuiltin("execute", functions.Execute(context.Background(), &fexec)),
			}

			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if !tt.hasErr && err != nil {
				t.Error(err)
			} else {
				return
			}

			if gotCmd != tt.wantCmd {
				t.Errorf("cmd: got: %s, want: %s", gotCmd, tt.wantCmd)
			}
			if diff := cmp.Diff(gotArgs, tt.wantArgs); diff != "" {
				t.Errorf("args diff: %s", diff)
			}
		})
	}
}
