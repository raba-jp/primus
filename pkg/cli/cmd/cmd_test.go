package cmd_test

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

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
		assert.NoErrorf(t, err, "Failed to read golden test file")
		assert.Equal(t, data, string(res))
	}
}
