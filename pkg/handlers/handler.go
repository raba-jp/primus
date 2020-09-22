//go:generate mockgen -destination mock/handler.go . CheckInstallHandler,InstallHandler,UninstallHandler,FileCopyHandler,FileMoveHandler,SymlinkHandler,HTTPRequestHandler,CommandHandler,FishSetVariableHandler,FishSetPathHandler,CreateDirectoryHandler

package handlers

import (
	"context"
	"os"
)

type CheckInstallHandler interface {
	CheckInstall(ctx context.Context, name string) bool
}

type CheckInstallHandlerFunc func(ctx context.Context, name string) bool

func (f CheckInstallHandlerFunc) CheckInstall(ctx context.Context, name string) bool {
	return f(ctx, name)
}

type InstallParams struct {
	Name   string
	Option string
}

type InstallHandler interface {
	Install(ctx context.Context, dryrun bool, p *InstallParams) error
}

type InstallHandlerFunc func(ctx context.Context, dryrun bool, p *InstallParams) error

func (f InstallHandlerFunc) Install(ctx context.Context, dryrun bool, p *InstallParams) error {
	return f(ctx, dryrun, p)
}

type UninstallParams struct {
	Name string
}

type UninstallHandler interface {
	Uninstall(ctx context.Context, dryrun bool, p *UninstallParams) error
}

type UninstallHandlerFunc func(ctx context.Context, dryrun bool, p *UninstallParams) error

func (f UninstallHandlerFunc) Uninstall(ctx context.Context, dryrun bool, p *UninstallParams) error {
	return f(ctx, dryrun, p)
}

type FileCopyParams struct {
	Src        string
	Dest       string
	Permission os.FileMode
}

type FileCopyHandler interface {
	FileCopy(ctx context.Context, dryrun bool, p *FileCopyParams) error
}

type FileCopyHandlerFunc func(ctx context.Context, dryrun bool, p *FileCopyParams) error

func (f FileCopyHandlerFunc) FileCopy(ctx context.Context, dryrun bool, p *FileCopyParams) error {
	return f(ctx, dryrun, p)
}

type FileMoveParams struct {
	Src  string
	Dest string
}

type FileMoveHandler interface {
	FileMove(ctx context.Context, dryrun bool, p *FileMoveParams) error
}

type FileMoveHandlerFunc func(ctx context.Context, dryrun bool, p *FileMoveParams) error

func (f FileMoveHandlerFunc) FileMove(ctx context.Context, dryrun bool, p *FileMoveParams) error {
	return f(ctx, dryrun, p)
}

type SymlinkParams struct {
	Src  string
	Dest string
	User string
}

type SymlinkHandler interface {
	Symlink(ctx context.Context, dryrun bool, p *SymlinkParams) error
}

type SymlinkHandlerFunc func(ctx context.Context, dryrun bool, p *SymlinkParams) error

func (f SymlinkHandlerFunc) Symlink(ctx context.Context, dryrun bool, p *SymlinkParams) error {
	return f(ctx, dryrun, p)
}

type HTTPRequestParams struct {
	URL  string
	Path string
}

type HTTPRequestHandler interface {
	HTTPRequest(ctx context.Context, dryrun bool, p *HTTPRequestParams) error
}

type HTTPRequestHandlerFunc func(ctx context.Context, dryrun bool, p *HTTPRequestParams) error

func (f HTTPRequestHandlerFunc) HTTPRequest(ctx context.Context, dryrun bool, p *HTTPRequestParams) error {
	return f(ctx, dryrun, p)
}

type CommandParams struct {
	CmdName string
	CmdArgs []string
	Cwd     string
	User    string
}

type CommandHandler interface {
	Command(ctx context.Context, dryrun bool, p *CommandParams) error
}

type CommandHandlerFunc func(ctx context.Context, dryrun bool, p *CommandParams) error

func (f CommandHandlerFunc) Command(ctx context.Context, dryrun bool, p *CommandParams) error {
	return f(ctx, dryrun, p)
}

type FishVariableScope int

const (
	FishVariableUniversalScope FishVariableScope = iota + 1
	FishVariableGlobalScope
	FishVariableLocalScope
)

type FishSetVariableParams struct {
	Name   string
	Value  string
	Scope  FishVariableScope
	Export bool
}

type FishSetVariableHandler interface {
	FishSetVariable(ctx context.Context, dryrun bool, p *FishSetVariableParams) error
}

type FishSetVariableHandlerFunc func(ctx context.Context, dryrun bool, p *FishSetVariableParams) error

func (f FishSetVariableHandlerFunc) FishSetVariable(ctx context.Context, dryrun bool, p *FishSetVariableParams) error {
	return f(ctx, dryrun, p)
}

type FishSetPathParams struct {
	Values []string
}

type FishSetPathHandler interface {
	FishSetPath(ctx context.Context, dryrun bool, p *FishSetPathParams) error
}

type FishSetPathHandlerFunc func(ctx context.Context, dryrun bool, p *FishSetPathParams) error

func (f FishSetPathHandlerFunc) FishSetPath(ctx context.Context, dryrun bool, p *FishSetPathParams) error {
	return f(ctx, dryrun, p)
}

type CreateDirectoryParams struct {
	Path       string
	Permission os.FileMode
}

type CreateDirectoryHandler interface {
	CreateDirectory(ctx context.Context, dryrun bool, p *CreateDirectoryParams) error
}

type CreateDirectoryHandlerFunc func(ctx context.Context, dryrun bool, p *CreateDirectoryParams) error

func (f CreateDirectoryHandlerFunc) FileExists(ctx context.Context, dryrun bool, p *CreateDirectoryParams) error {
	return f(ctx, dryrun, p)
}
