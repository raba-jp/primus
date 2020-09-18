package backend

import (
	"time"

	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/raba-jp/primus/pkg/internal/handlers"
	"github.com/spf13/afero"
)

const installTimeout = 5 * time.Minute

type Backend interface {
	handlers.CheckInstallHandler
	handlers.InstallHandler
	handlers.UninstallHandler
	handlers.FileCopyHandler
	handlers.FileMoveHandler
	handlers.SymlinkHandler
	handlers.HTTPRequestHandler
	handlers.CommandHandler
	handlers.FishSetVariableHandler
	handlers.FishSetPathHandler
}

func NewFs() afero.Fs {
	return afero.NewOsFs()
}

func NewExecInterface() exec.Interface {
	return exec.New()
}

func New(execIF exec.Interface, fs afero.Fs) Backend {
	os := DetectOS(execIF, fs)
	switch os {
	case Darwin:
		return &DarwinBackend{
			Exec: execIF,
			Fs:   fs,
			BaseBackend: &BaseBackend{
				Exec: execIF,
				Fs:   fs,
			},
		}
	case Arch:
		return &ArchLinuxBackend{
			Exec: execIF,
			BaseBackend: &BaseBackend{
				Exec: execIF,
				Fs:   fs,
			},
		}
	case Unknown:
		return nil
	}
	return nil
}
