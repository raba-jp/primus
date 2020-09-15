package cmd_test

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

func executeCommand(child *cobra.Command, writer io.Writer, args ...string) error {
	p := cobra.Command{}
	p.AddCommand(child)
	p.SetOut(writer)
	p.SetErr(writer)

	a := append([]string{}, child.Use)
	a = append(a, args...)
	p.SetArgs(a)

	return p.Execute()
}

func goldenTest(t *testing.T, path string, data string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		ioutil.WriteFile(path, []byte(data), 0644)
		return
	} else {
		res, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read golden test file: %s: %v", path, err)
		}
		if diff := cmp.Diff(string(res), data); diff != "" {
			t.Fatal(diff)
		}
	}
}
