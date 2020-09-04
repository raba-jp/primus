package apply_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/exec"
	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"
	"github.com/raba-jp/primus/pkg/executor"
)

func TestExecutor_Command(t *testing.T) {
	tests := []struct {
		name       string
		cmd        string
		args       []string
		mockStdout string
	}{
		{
			name:       "echo hello world",
			cmd:        "echo",
			args:       []string{"hello", "world"},
			mockStdout: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Reset()

			execIF = &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						fake := &fakeexec.FakeCmd{
							// Stdout: new(bytes.Buffer),
							RunScript: []fakeexec.FakeAction{
								func() ([]byte, []byte, error) {
									return []byte(tt.mockStdout), []byte{}, nil
								},
							},
						}
						return fakeexec.InitFakeCmd(fake, cmd, args...)
					},
				},
			}

			exc := NewTestExecutor()
			ret, err := exc.Command(context.Background(), &executor.CommandParams{
				CmdName: tt.cmd,
				CmdArgs: tt.args,
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
			if !ret {
				t.Fatalf("Failed to exec command: %s %s", tt.cmd, tt.args)
			}
		})
	}
}
