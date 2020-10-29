//go:generate mockery -outpkg=mocks -case=snake -name=ExecutableHandler

package handlers

import (
	"context"
	"os"
	"strings"

	"github.com/raba-jp/primus/pkg/env"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

type shell int

const (
	posixShell shell = iota + 1
	fishShell
)

type ExecutableHandler interface {
	Run(ctx context.Context, name string) (ok bool)
}

type ExecutableHandlerFunc func(ctx context.Context, name string) bool

func (f ExecutableHandlerFunc) Run(ctx context.Context, name string) bool {
	return f(ctx, name)
}

func NewExecutable(fs afero.Fs) ExecutableHandler {
	return ExecutableHandlerFunc(func(ctx context.Context, name string) bool {
		path := env.Get("PATH")

		paths := strings.Split(path, ":")
		executable := false
		walkFn := func(path string, info os.FileInfo, err error) error {
			executable = executable || strings.HasSuffix(path, name)
			return nil
		}
		for _, p := range paths {
			_ = afero.Walk(fs, p, walkFn)
		}

		zap.L().Debug(
			"check executable",
			zap.String("name", name),
			zap.Strings("path", paths),
			zap.Bool("ok", executable),
		)

		return executable
	})
}
