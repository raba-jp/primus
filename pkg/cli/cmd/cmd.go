package cmd

import (
	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/spf13/afero"
)

func NewFs() afero.Fs {
	return afero.NewOsFs()
}

func NewExecInterface() exec.Interface {
	return exec.New()
}
