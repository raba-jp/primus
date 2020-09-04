package cmd_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/cmd"
	"github.com/raba-jp/primus/pkg/cli/ui"
)

func TestVersion(t *testing.T) {
	buf := new(bytes.Buffer)
	ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

	versionCmd := cmd.NewVersionCommand()
	if err := executeCommand(versionCmd); err != nil {
		t.Fatalf("%v", err)
	}
	if diff := cmp.Diff("unset", buf.String()); diff != "" {
		t.Fatal(diff)
	}
}
