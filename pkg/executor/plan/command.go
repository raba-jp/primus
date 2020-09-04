package plan

import (
	"bytes"
	"context"
	"fmt"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
)

func (e *planExecutor) Command(ctx context.Context, p *executor.CommandParams) (bool, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s ", p.CmdName)
	for _, arg := range p.CmdArgs {
		fmt.Fprintf(buf, "%s ", arg)
	}
	fmt.Fprintf(buf, "\n")

	ui.Printf(buf.String())
	return true, nil
}
