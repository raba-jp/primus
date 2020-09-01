package plan

import (
	"io"

	"github.com/raba-jp/primus/executor"
)

type planExecutor struct {
	Out    io.Writer
	Errout io.Writer
}

func NewPlanExecutorWithArgs(out io.Writer, errout io.Writer) executor.Executor {
	return &planExecutor{
		Out:    out,
		Errout: errout,
	}
}
