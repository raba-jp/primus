package modules

import (
	"github.com/raba-jp/primus/pkg/exec"
)

func NewExecInterface() exec.Interface {
	return exec.New()
}
