package backend

import (
	"net/http"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/spf13/afero"
)

func NewFs() afero.Fs {
	return afero.NewOsFs()
}

func NewExecInterface() exec.Interface {
	return exec.New()
}

func NewHTTPClient() *http.Client {
	return http.DefaultClient
}
