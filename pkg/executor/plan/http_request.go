package plan

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
)

func (e *planExecutor) HTTPRequest(ctx context.Context, p *executor.HTTPRequestParams) (bool, error) {
	ui.Printf("curl -Lo %s %s\n", p.Path, p.URL)
	return true, nil
}
