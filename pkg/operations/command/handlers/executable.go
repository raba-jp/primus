//go:generate mockgen -destination mock/executable.go . ExecutableHandler

package handlers

import (
	"context"
	"os"
	"strings"

	"github.com/raba-jp/primus/pkg/env"
	"github.com/spf13/afero"
)

type shell int

const (
	posixShell shell = iota + 1
	fishShell
)

type ExecutableHandler interface {
	Executable(ctx context.Context, name string) (ok bool)
}

type ExecutableHandlerFunc func(ctx context.Context, name string) (ok bool)

func (f ExecutableHandlerFunc) Executable(ctx context.Context, name string) (ok bool) {
	return f(ctx, name)
}

func NewExecutable(fs afero.Fs) ExecutableHandler {
	return ExecutableHandlerFunc(func(ctx context.Context, name string) bool {
		path := env.Get("PATH")
		separator := ":"
		if shell := detectShell(); shell == fishShell {
			separator = " "
		}

		paths := strings.Split(path, separator)
		executable := false
		walkFn := func(path string, info os.FileInfo, err error) error {
			executable = executable || strings.HasSuffix(path, name)
			return nil
		}
		for _, p := range paths {
			_ = afero.Walk(fs, p, walkFn)
		}

		return executable
	})
}

func detectShell() shell {
	shell := env.Get("SHELL")
	if strings.HasSuffix(shell, "fish") {
		return fishShell
	}
	return posixShell
}
