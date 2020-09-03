package apply

import (
	"net/http"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/spf13/afero"
)

type applyExecutor struct {
	Exec   exec.Interface
	Fs     afero.Fs
	Client *http.Client
}

func NewApplyExecutror() executor.Executor {
	return &applyExecutor{
		Exec:   exec.New(),
		Fs:     afero.NewOsFs(),
		Client: http.DefaultClient,
	}
}

func NewApplyExecutorWithArgs(exc exec.Interface, fs afero.Fs, client *http.Client) executor.Executor {
	return &applyExecutor{
		Exec:   exc,
		Fs:     fs,
		Client: client,
	}
}
