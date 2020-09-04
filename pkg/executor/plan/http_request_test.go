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

func TestHTTPRequest(t *testing.T) {
	tests := []struct {
		name string
		url  string
		path string
		want string
	}{
		{
			name: "success",
			url:  "https://example.com",
			path: "/sym/output.txt",
			want: "curl -Lo /sym/output.txt https://example.com\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

			exc := plan.NewPlanExecutor()
			ok, err := exc.HTTPRequest(context.Background(), &executor.HTTPRequestParams{
				URL:  tt.url,
				Path: tt.path,
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
