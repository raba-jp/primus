package plan

import (
	"context"
	"fmt"

	"github.com/raba-jp/primus/executor"
)

func (e *planExecutor) FileMove(ctx context.Context, p *executor.FileMoveParams) (bool, error) {
	fmt.Fprintf(e.Out, "%s => %s", p.Src, p.Dest)
	return true, nil
}
