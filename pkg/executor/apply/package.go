package apply

import (
	"context"

	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/internal/writer"
	"golang.org/x/xerrors"
)

func (e *applyExecutor) Package(ctx context.Context, p *executor.PackageParams) (bool, error) {
	cmd := e.Exec.CommandContext(ctx, "pacman", "-S", "--noconfirm", p.Name)
	nop := writer.NopWriter{}
	cmd.SetStdout(&nop)
	cmd.SetStderr(&nop)
	if err := cmd.Run(); err != nil {
		return false, xerrors.Errorf(": %w", err)
	}

	return true, nil
}
