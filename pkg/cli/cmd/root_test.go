package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/raba-jp/primus/pkg/cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		goldenFile string
	}{
		{
			name:       "no args",
			args:       []string{},
			goldenFile: "execute_no_args.golden",
		},
		{
			name:       "set logLevel",
			args:       []string{"--logLevel=info"},
			goldenFile: "execute_enable_debug.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buf := new(bytes.Buffer)
			rootCmd := cmd.Initialize()
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()
			assert.NoError(t, err)

			wd, _ := os.Getwd()
			path := filepath.Join(wd, "testdata", "golden", tt.goldenFile)
			goldenTest(t, path, buf.String())
		})
	}
}
