package apply

import (
	"io"
	"net/http"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/spf13/afero"
)

type applyExecutor struct {
	In     io.Reader
	Out    io.Writer
	Errout io.Writer
	Exec   exec.Interface
	Fs     afero.Fs
	Client *http.Client
}

func NewApplyExecutorWithArgs(in io.Reader, out io.Writer, errout io.Writer, exc exec.Interface, fs afero.Fs, client *http.Client) executor.Executor {
	return &applyExecutor{
		In:     in,
		Out:    out,
		Errout: errout,
		Exec:   exc,
		Fs:     fs,
		Client: client,
	}
}
