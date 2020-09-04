package plan

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
)

func (e *planExecutor) Package(ctx context.Context, p *executor.PackageParams) (bool, error) {
	ui.Printf("%s\n", p.Name)
	return true, nil
}
