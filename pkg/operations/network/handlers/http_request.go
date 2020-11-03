//go:generate mockery -outpkg=mocks -case=snake -name=HTTPRequestHandler

package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

const timeout = 10 * time.Minute

type HTTPRequestParams struct {
	URL  string
	Path string
}

type HTTPRequestHandler interface {
	Run(ctx context.Context, p *HTTPRequestParams) (err error)
}

type HTTPRequestHandlerFunc func(ctx context.Context, p *HTTPRequestParams) error

func (f HTTPRequestHandlerFunc) Run(ctx context.Context, p *HTTPRequestParams) error {
	return f(ctx, p)
}

func NewHTTPRequest(client *http.Client, fs afero.Fs) HTTPRequestHandler {
	return HTTPRequestHandlerFunc(func(ctx context.Context, p *HTTPRequestParams) error {
		if dryrun := ctxlib.DryRun(ctx); dryrun {
			httpRequestDryRun(p)
			return nil
		}
		return httpRequest(ctx, client, fs, p)
	})
}

func httpRequest(ctx context.Context, client *http.Client, fs afero.Fs, p *HTTPRequestParams) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, p.URL, nil)
	if err != nil {
		return xerrors.Errorf("Failed to create new http request: %w", err)
	}
	req = req.WithContext(ctx)

	res, err := client.Do(req)
	if err != nil {
		return xerrors.Errorf("Failed to http request: %w", err)
	}
	defer res.Body.Close()

	if err := afero.WriteReader(fs, p.Path, res.Body); err != nil {
		return xerrors.Errorf("Failed to write response body: %w", err)
	}

	return nil
}

func httpRequestDryRun(p *HTTPRequestParams) {
	ui.Printf("curl -Lo %s %s\n", p.Path, p.URL)
}
