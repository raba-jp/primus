package plan

import (
	"io"

	"github.com/raba-jp/primus/pkg/executor"
)

type planExecutor struct {
	Out io.Writer
}

func NewPlanExecutorWithArgs(out io.Writer) executor.Executor {
	return &planExecutor{
		Out: out,
	}
}
