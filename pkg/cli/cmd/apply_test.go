package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/raba-jp/primus/pkg/cli/cmd"
	"github.com/raba-jp/primus/pkg/cli/ui"
)

func TestApply(t *testing.T) {
	wd, _ := os.Getwd()

	tests := []struct {
		name       string
		args       []string
		hasErr     bool
		goldenFile string
	}{
		{
			name:       "no args",
			args:       []string{},
			hasErr:     true,
			goldenFile: "apply_no_args.golden",
		},
		{
			name: "success",
			args: []string{
				filepath.Join(wd, "testdata", "apply.star"),
			},
			hasErr:     false,
			goldenFile: "apply_success.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})
			applyCmd := cmd.NewApplyCommand()
			if err := executeCommand(applyCmd, buf, tt.args...); !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
			path := filepath.Join(wd, "testdata", "golden", tt.goldenFile)
			goldenTest(t, path, buf.String())
		})
	}
}
