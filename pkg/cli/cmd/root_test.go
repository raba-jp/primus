package cmd_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/cmd"
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
			rootCmd := cmd.NewPrimusCommand()
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("%v", err)
			}

			wd, _ := os.Getwd()
			path := filepath.Join(wd, "testdata", "golden", tt.goldenFile)
			if _, err := os.Stat(path); err != nil {
				ioutil.WriteFile(path, buf.Bytes(), 0644)
				return
			} else {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					t.Fatalf("%v", err)
				}
				if diff := cmp.Diff(data, buf.Bytes()); diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}
