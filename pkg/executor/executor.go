//go:generate mockgen -destination mock/executor.go . Executor

package executor

import (
	"context"
	"os"
)

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

type HTTPRequestParams struct {
	URL  string
	Path string
}

type PackageParams struct {
	Name   string
	Option string
}

type SymlinkParams struct {
	Src  string
	Dest string
	User string
}

type Executor interface {
	Command(ctx context.Context, p *CommandParams) (bool, error)
	Symlink(ctx context.Context, p *SymlinkParams) (bool, error)
	FileCopy(ctx context.Context, p *FileCopyParams) (bool, error)
	FileMove(ctx context.Context, p *FileMoveParams) (bool, error)
	HTTPRequest(ctx context.Context, p *HTTPRequestParams) (bool, error)
	Package(ctx context.Context, p *PackageParams) (bool, error)
}
