package plan

import (
	"context"
	"fmt"

	"github.com/raba-jp/primus/executor"
)

func (e *planExecutor) Symlink(ctx context.Context, p *executor.SymlinkParams) (bool, error) {
	fmt.Fprintf(e.Out, "%s => %s", p.Src, p.Dest)
	return true, nil
}
