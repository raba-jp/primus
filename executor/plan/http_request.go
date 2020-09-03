package plan

import (
	"context"
	"fmt"

	"github.com/raba-jp/primus/executor"
)

func (e *planExecutor) HTTPRequest(ctx context.Context, p *executor.HTTPRequestParams) (bool, error) {
	fmt.Fprintf(e.Out, "URL: %s SavePath: %s", p.URL, p.Path)
	return true, nil
}
