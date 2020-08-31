package executor

import (
	"context"

	"golang.org/x/xerrors"
)

type PackageParams struct {
	Name string
}

func (e *executor) Package(ctx context.Context, p *PackageParams) (bool, error) {
	cmd := e.Exec.CommandContext(ctx, "pacman", "-S", "--noconfirm", p.Name)
	cmd.SetStdout(e.Out)
	cmd.SetStderr(e.Errout)
	if err := cmd.Run(); err != nil {
		return false, xerrors.Errorf(": %w", err)
	}

	return true, nil
}
