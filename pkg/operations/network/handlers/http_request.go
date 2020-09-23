//go:generate mockgen -destination mock/handler.go . HTTPRequestHandler

package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/k0kubun/pp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

const timeout = 10 * time.Minute

func init() {
	pp.ColoringEnabled = false
}

type HTTPRequestParams struct {
	URL  string
	Path string
}

func (p *HTTPRequestParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type HTTPRequestHandler interface {
	HTTPRequest(ctx context.Context, dryrun bool, p *HTTPRequestParams) error
}

type HTTPRequestHandlerFunc func(ctx context.Context, dryrun bool, p *HTTPRequestParams) error

func (f HTTPRequestHandlerFunc) HTTPRequest(ctx context.Context, dryrun bool, p *HTTPRequestParams) error {
	return f(ctx, dryrun, p)
}

func NewHTTPRequest(client *http.Client, fs afero.Fs) HTTPRequestHandler {
	return HTTPRequestHandlerFunc(func(ctx context.Context, dryrun bool, p *HTTPRequestParams) error {
		if dryrun {
			ui.Printf("curl -Lo %s %s\n", p.Path, p.URL)
			return nil
		}

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
	})
}
