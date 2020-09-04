package plan

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
)

func (e *planExecutor) Symlink(ctx context.Context, p *executor.SymlinkParams) (bool, error) {
	ui.Printf("ln -s %s %s\n", p.Src, p.Dest)
	return true, nil
}
