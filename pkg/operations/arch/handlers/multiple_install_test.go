package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/operations/arch/handlers"
	"github.com/raba-jp/primus/pkg/operations/command/handlers/mocks"
)

func TestNewMultipleInstall(t *testing.T) {
	executable := mocks.ExecutableHandlerRunExpectation{
		Args: mocks.ExecutableHandlerRunArgs{
			CtxAnything:  true,
			NameAnything: true,
		},
		Returns: mocks.ExecutableHandlerRunReturns{Ok: true},
	}
	tests := []struct {
		name       string
		executable mocks.ExecutableHandlerRunExpectation
		mock       exec.InterfaceCommandContextExpectation
		errAssert  assert.ErrorAssertionFunc
	}{
		{
			name:       "success",
			executable: executable,
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "yay",
					Args:        []string{"--pacman", "powerpill", "-S", "--noconfirm", "arg1", "arg2"},
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyRunExpectation(exec.CmdRunExpectation{
							Returns: exec.CmdRunReturns{Err: nil},
						})
						return cmd
					},
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:       "failed",
			executable: executable,
			mock: exec.InterfaceCommandContextExpectation{
				Args: exec.InterfaceCommandContextArgs{
					CtxAnything: true,
					Cmd:         "yay",
					Args:        []string{"--pacman", "powerpill", "-S", "--noconfirm", "arg1", "arg2"},
				},
				Returns: exec.InterfaceCommandContextReturns{
					Cmd: func() exec.Cmd {
						cmd := new(exec.MockCmd)
						cmd.ApplyRunExpectation(exec.CmdRunExpectation{
							Returns: exec.CmdRunReturns{Err: xerrors.New("dummy")},
						})
						return cmd
					},
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		executable := new(mocks.ExecutableHandler)
		executable.ApplyRunExpectation(tt.executable)

		exc := new(exec.MockInterface)
		exc.ApplyCommandContextExpectation(tt.mock)

		multipleInstall := handlers.NewMultipleInstall(executable, exc)
		err := multipleInstall.Run(context.Background(), false, &handlers.MultipleInstallParams{
			Names: []string{"arg1", "arg2"},
		})
		tt.errAssert(t, err)
	}
}

func TestNewMultipleInstall__dryrun(t *testing.T) {
	buf := new(bytes.Buffer)
	ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})

	executable := new(mocks.ExecutableHandler)
	executable.ApplyRunExpectation(mocks.ExecutableHandlerRunExpectation{
		Args: mocks.ExecutableHandlerRunArgs{
			CtxAnything:  true,
			NameAnything: true,
		},
		Returns: mocks.ExecutableHandlerRunReturns{Ok: true},
	})

	multipleInstall := handlers.NewMultipleInstall(executable, nil)
	err := multipleInstall.Run(context.Background(), true, &handlers.MultipleInstallParams{
		Names: []string{"arg1", "arg2"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "yay --pacman powerpill -S --noconfirm arg1 arg2\n", buf.String())
}
