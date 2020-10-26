//go:generate mockery -outpkg=mocks -case=snake -name=CheckInstallHandler

package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/spf13/afero"
)

type CheckInstallHandler interface {
	Run(ctx context.Context, name string) (ok bool)
}

type CheckInstallHandlerFunc func(ctx context.Context, name string) bool

func (f CheckInstallHandlerFunc) Run(ctx context.Context, name string) bool {
	return f(ctx, name)
}

func NewCheckInstall(exc exec.Interface, fs afero.Fs) CheckInstallHandler {
	return CheckInstallHandlerFunc(func(ctx context.Context, name string) bool {
		installed := false
		walkFn := func(path string, info os.FileInfo, err error) error {
			installed = installed || strings.Contains(path, name)
			return nil
		}

		// brew list
		res, _ := exc.CommandContext(ctx, "brew", "--prefix").Output()
		prefix := strings.ReplaceAll(string(res), "\n", "")
		_ = afero.Walk(fs, fmt.Sprintf("%s/Celler", prefix), walkFn)

		// brew cask list
		_ = afero.Walk(fs, "/opt/homebrew-cask/Caskroom", walkFn)
		_ = afero.Walk(fs, "/usr/local/Caskroom", walkFn)

		return installed
	})
}
