package apply

import (
	"context"

	"github.com/raba-jp/primus/executor"
	"golang.org/x/xerrors"
)

func (e *applyExecutor) Package(ctx context.Context, p *executor.PackageParams) (bool, error) {
	cmd := e.Exec.CommandContext(ctx, "pacman", "-S", "--noconfirm", p.Name)
	cmd.SetStdout(e.Out)
	cmd.SetStderr(e.Errout)
	if err := cmd.Run(); err != nil {
		return false, xerrors.Errorf(": %w", err)
	}

	return true, nil
}
