package apply

import (
	"context"
	"net/http"
	"time"

	"github.com/raba-jp/primus/executor"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

func (e *applyExecutor) HTTPRequest(ctx context.Context, p *executor.HTTPRequestParams) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, p.URL, nil)
	if err != nil {
		return false, xerrors.Errorf("Failed to create new http request: %w", err)
	}
	req = req.WithContext(ctx)

	res, err := e.Client.Do(req)
	if err != nil {
		return false, xerrors.Errorf("Failed to http request: %w", err)
	}
	defer res.Body.Close()

	if err := afero.WriteReader(e.Fs, p.Path, res.Body); err != nil {
		return false, xerrors.Errorf("Failed to write response body: %w", err)
	}

	return true, nil
}
