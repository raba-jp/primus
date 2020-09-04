package plan_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/executor/plan"
)

func TestFileMove(t *testing.T) {
	tests := []struct {
		name string
		src  string
		dest string
		want string
	}{
		{
			name: "success",
			src:  "/sym/src.txt",
			dest: "/sym/dest.txt",
			want: "mv /sym/src.txt /sym/dest.txt\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			exc := plan.NewPlanExecutor()
			ok, err := exc.FileMove(context.Background(), &executor.FileMoveParams{
				Src:  tt.src,
				Dest: tt.dest,
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
			if !ok {
				t.Fatal("Unexpected error")
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
