package executor

import (
	"context"

	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

type HttpRequestParams struct {
	URL  string
	Path string
}

func (e *executor) HttpRequest(ctx context.Context, p *HttpRequestParams) (bool, error) {
	res, err := e.Client.Get(p.URL)
	if err != nil {
		return false, xerrors.Errorf("Failed to http request: %w", err)
	}
	defer res.Body.Close()

	if err := afero.WriteReader(e.Fs, p.Path, res.Body); err != nil {
		return false, xerrors.Errorf("Failed to write response body: %w", err)
	}

	return true, nil
}
