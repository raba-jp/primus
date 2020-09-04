package plan

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
)

func (e *planExecutor) FileMove(ctx context.Context, p *executor.FileMoveParams) (bool, error) {
	ui.Printf("mv %s %s\n", p.Src, p.Dest)
	return true, nil
}
