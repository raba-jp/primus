package plan

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
)

func TestCommand(t *testing.T) {
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

			exc := NewPlanExecutor()
			ok, err := exc.Command(context.Background(), &executor.CommandParams{
				CmdName: tt.command,
				CmdArgs: tt.args,
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
			if !ok {
				t.Fatal("Unexpected error")
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
