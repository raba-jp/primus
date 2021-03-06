package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/cli/cmd"
	"github.com/raba-jp/primus/pkg/cli/ui"
)

func TestApply(t *testing.T) {
	wd, _ := os.Getwd()

	tests := []struct {
		name       string
		args       []string
		errAssert  assert.ErrorAssertionFunc
		goldenFile string
	}{
		{
			name:       "no args",
			args:       []string{},
			errAssert:  assert.Error,
			goldenFile: "apply_no_args.golden",
		},
		{
			name: "success",
			args: []string{
				filepath.Join(wd, "testdata", "apply.star"),
			},
			errAssert:  assert.NoError,
			goldenFile: "apply_success.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})
			applyCmd := cmd.NewApplyCommand()
			err := executeCommand(applyCmd, buf, tt.args...)
			tt.errAssert(t, err)
			path := filepath.Join(wd, "testdata", "golden", tt.goldenFile)
			goldenTest(t, path, buf.String())
		})
	}
}
