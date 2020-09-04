package plan

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
)

func (e *planExecutor) FileCopy(ctx context.Context, p *executor.FileCopyParams) (bool, error) {
	ui.Printf("cp %s %s\n", p.Src, p.Dest)
	return true, nil
}
