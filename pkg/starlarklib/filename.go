package starlarklib

import (
	"go.starlark.net/starlark"
	"go.uber.org/zap"
)

func GetCurrentFilePath(thread *starlark.Thread) string {
	callframe := thread.CallFrame(thread.CallStackDepth() - 1)
	path := callframe.Pos.Filename()
	zap.L().Debug(
		"callframe",
		zap.Int("depth", thread.CallStackDepth()),
		zap.String("filepath", path),
	)
	return path
}
