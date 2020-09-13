//go:generate mockgen -destination mock/backend.go . Backend
package backend

import (
	"context"
	"os"
	"time"

	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/spf13/afero"
)

const installTimeout = 5 * time.Minute

type InstallParams struct {
	Name   string
	Option string
}

type CommandParams struct {
	CmdName string
	CmdArgs []string
	Cwd     string
	User    string
}

type FileCopyParams struct {
	Src        string
	Dest       string
	Permission os.FileMode
}

type FileMoveParams struct {
	Src  string
	Dest string
}

type SymlinkParams struct {
	Src  string
	Dest string
	User string
}

type HTTPRequestParams struct {
	URL  string
	Path string
}

type Backend interface {
	CheckInstall(ctx context.Context, name string) bool
	Install(ctx context.Context, p *InstallParams) error
	Uninstall(ctx context.Context, name string) error
	FileCopy(ctx context.Context, p *FileCopyParams) error
	FileMove(ctx context.Context, p *FileMoveParams) error
	Symlink(ctx context.Context, p *SymlinkParams) error
	HTTPRequest(ctx context.Context, p *HTTPRequestParams) error
	Command(ctx context.Context, p *CommandParams) error
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
