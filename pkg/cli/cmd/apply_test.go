package cmd_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/cli/cmd"
)

func TestApply(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		hasErr bool
	}{
		{
			name:   "no args",
			args:   []string{},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyCmd := cmd.NewApplyCommand()
			if err := executeCommand(applyCmd, tt.args...); !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
