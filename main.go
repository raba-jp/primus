package main

import (
	"github.com/raba-jp/starlark_iac/cmd"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	cmd.Execute()
}
