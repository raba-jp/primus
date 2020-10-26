//go:generate mockery -outpkg=mocks -case=snake -name=CheckInstallHandler

package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/exec"
)

type CheckInstallHandler interface {
	Run(ctx context.Context, name string) (ok bool)
}

type CheckInstallHandlerFunc func(ctx context.Context, name string) bool

func (f CheckInstallHandlerFunc) Run(ctx context.Context, name string) bool {
	return f(ctx, name)
}

func NewCheckInstall(exc exec.Interface) CheckInstallHandler {
	return CheckInstallHandlerFunc(func(ctx context.Context, name string) bool {
		err := exc.CommandContext(ctx, "pacman", "-Qg", name).Run()
		return err == nil
	})
}
