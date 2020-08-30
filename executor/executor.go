//go:generate mockgen -destination mock/executor.go . Executor

package executor

import (
	"context"
	"io"
	"net/http"

	"github.com/raba-jp/primus/exec"
	"github.com/spf13/afero"
)

type Executor interface {
	Command(ctx context.Context, p *CommandParams) (bool, error)
	Symlink(ctx context.Context, p *SymlinkParams) (bool, error)
	FileCopy(ctx context.Context, p *FileCopyParams) (bool, error)
	FileMove(ctx context.Context, p *FileMoveParams) (bool, error)
	HttpRequest(ctx context.Context, p *HttpRequestParams) (bool, error)
	Package(ctx context.Context, p *PackageParams) (bool, error)
}

type executor struct {
	In     io.Reader
	Out    io.Writer
	Errout io.Writer
	Exec   exec.Interface
	Fs     afero.Fs
	Client *http.Client
}

func NewExecutorWithArgs(in io.Reader, out io.Writer, errout io.Writer, exc exec.Interface, fs afero.Fs, client *http.Client) Executor {
	return &executor{
		In:     in,
		Out:    out,
		Errout: errout,
		Exec:   exc,
		Fs:     fs,
		Client: client,
	}
}
