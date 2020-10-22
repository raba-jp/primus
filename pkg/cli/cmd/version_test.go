package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/cli/cmd"
	"github.com/raba-jp/primus/pkg/cli/ui"
)

func TestVersion(t *testing.T) {
	buf := new(bytes.Buffer)
	ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

	versionCmd := cmd.NewVersionCommand()
	err := executeCommand(versionCmd, buf)
	assert.NoError(t, err)
	assert.Equal(t, "unset", buf.String())
}
