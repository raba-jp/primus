//go:generate mockgen -destination mock/handler.go . CheckInstallHandler,InstallHandler,UninstallHandler

package handlers

import (
	"context"
	"time"

	"github.com/k0kubun/pp"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/spf13/afero"
)

const installTimeout = 5 * time.Minute

func init() {
	pp.ColoringEnabled = false
}

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

func (p *InstallParams) String() string {
	return pp.Sprintf("%v\n", p)
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

func (p *UninstallParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type UninstallHandler interface {
	Uninstall(ctx context.Context, dryrun bool, p *UninstallParams) error
}

type UninstallHandlerFunc func(ctx context.Context, dryrun bool, p *UninstallParams) error

func (f UninstallHandlerFunc) Uninstall(ctx context.Context, dryrun bool, p *UninstallParams) error {
	return f(ctx, dryrun, p)
}

type osHandler interface {
	CheckInstallHandler
	InstallHandler
	UninstallHandler
}

func NewCheckInstall(execIF exec.Interface, fs afero.Fs) CheckInstallHandler {
	h := getOSHandler(execIF, fs)
	return CheckInstallHandlerFunc(h.CheckInstall)
}

func NewInstall(execIF exec.Interface, fs afero.Fs) InstallHandler {
	h := getOSHandler(execIF, fs)
	return InstallHandlerFunc(h.Install)
}

func NewUninstall(execIF exec.Interface, fs afero.Fs) UninstallHandler {
	h := getOSHandler(execIF, fs)
	return UninstallHandlerFunc(h.Uninstall)
}

func getOSHandler(execIF exec.Interface, fs afero.Fs) osHandler {
	switch backend.DetectOS(execIF, fs) {
	case backend.Arch:
		return &archLinux{Exec: execIF}
	case backend.Darwin:
		return &darwin{Exec: execIF, Fs: fs}
	case backend.Unknown:
		fallthrough
	default:
		return nil
	}
}
