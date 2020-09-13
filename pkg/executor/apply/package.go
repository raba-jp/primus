package apply

import (
	"context"

	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/internal/backend"
	"golang.org/x/xerrors"
)

func (e *applyExecutor) Package(ctx context.Context, p *executor.PackageParams) (bool, error) {
	be := backend.New(e.Exec, e.Fs)
	if chk := be.CheckInstall(ctx, p.Name); chk {
		return true, nil
	}
	if err := be.Install(ctx, p.Name, p.Option); err != nil {
		return false, xerrors.Errorf(": %w", err)
	}
	return true, nil
}
