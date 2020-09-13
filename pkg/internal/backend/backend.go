package backend

import (
	"context"
	"time"

	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/spf13/afero"
)

const installTimeout = 5 * time.Minute

type Backend interface {
	CheckInstall(ctx context.Context, name string) bool
	Install(ctx context.Context, name string, option string) error
	Uninstall(ctx context.Context, name string) error
}

func New(execIF exec.Interface, fs afero.Fs) Backend {
	os := DetectOS(execIF, fs)
	switch os {
	case Darwin:
		return &DarwinBackend{
			Exec: execIF,
			Fs:   fs,
		}
	case Manjaro:
		return &ArchLinuxBackend{
			Exec: execIF,
		}
	case Unknown:
		return nil
	}
	return nil
}
