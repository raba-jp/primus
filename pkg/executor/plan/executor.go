package plan

import (
	"github.com/raba-jp/primus/pkg/executor"
)

type planExecutor struct {
}

func NewPlanExecutor() executor.Executor {
	return &planExecutor{}
}
