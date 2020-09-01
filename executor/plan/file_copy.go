package plan

import (
	"context"
	"fmt"

	"github.com/raba-jp/primus/executor"
)

func (e *planExecutor) FileCopy(ctx context.Context, p *executor.FileCopyParams) (bool, error) {
	fmt.Fprintf(e.Out, "%s => %s", p.Src, p.Dest)
	return true, nil
}
