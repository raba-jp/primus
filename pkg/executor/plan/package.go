package plan

import (
	"context"
	"fmt"

	"github.com/raba-jp/primus/pkg/executor"
)

func (e *planExecutor) Package(ctx context.Context, p *executor.PackageParams) (bool, error) {
	fmt.Fprintf(e.Out, "%s", p.Name)
	return true, nil
}
