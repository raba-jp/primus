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

func TestPackage(t *testing.T) {
	tests := []struct {
		data     string
		wantCmd  string
		wantArgs []string
		hasErr   bool
	}{
		{
			data:     `package(name="base-devel")`,
			wantCmd:  "pacman",
			wantArgs: []string{"-S", "--noconfirm", "base-devel"},
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
				"package": starlark.NewBuiltin("package", functions.Execute(context.Background(), &fexec)),
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
