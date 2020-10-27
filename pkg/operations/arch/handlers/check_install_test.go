package handlers_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/arch/handlers"
	"github.com/stretchr/testify/assert"
)

func TestNewCheckInstall(t *testing.T) {
	tests := []struct {
		name string
		mock exec.InterfaceCommandContextExpectation
		want bool
	}{
		{
			name: "success",
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "pacman",
					Args:        []string{"-Qg", "base-devel"},
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyRunExpectation(exec.CmdRunExpectation{
							Returns: exec.CmdRunReturns{
								Err: nil,
							},
						})
						return cmd
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exc := new(exec.MockInterface)
			exc.ApplyCommandContextExpectation(tt.mock)

			checkInstall := handlers.NewCheckInstall(exc)
			res := checkInstall.Run(context.Background(), "base-devel")
			assert.Equal(t, tt.want, res)
		})
	}
}
