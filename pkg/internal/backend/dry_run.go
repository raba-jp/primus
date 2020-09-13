package backend

import (
	"bytes"
	"context"
	"fmt"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/spf13/afero"
)

type DryRunBackend struct {
	Exec exec.Interface
	Fs   afero.Fs
}

func NewDryRunBackend(execIF exec.Interface, fs afero.Fs) Backend {
	return &DryRunBackend{
		Exec: execIF,
		Fs:   fs,
	}
}

func (b *DryRunBackend) CheckInstall(ctx context.Context, name string) bool {
	return false
}

func (b *DryRunBackend) Install(ctx context.Context, p *InstallParams) error {
	switch DetectOS(b.Exec, b.Fs) {
	case Darwin:
		ui.Printf("brew install %s %s", p.Option, p.Name)
	case Arch:
		ui.Printf("pacman -S --noconfirm %s %s", p.Option, p.Name)
	case Unknown:
		ui.Printf("Unsupported OS")
	}
	return nil
}

func (b *DryRunBackend) Uninstall(ctx context.Context, name string) error {
	switch DetectOS(b.Exec, b.Fs) {
	case Darwin:
		ui.Printf("brew uninstall %s", name)
	case Arch:
		ui.Printf("pacman -R %s", name)
	case Unknown:
		ui.Printf("Unsupported OS")
	}
	return nil
}

func (b *DryRunBackend) FileCopy(ctx context.Context, p *FileCopyParams) error {
	ui.Printf("cp %s %s\n", p.Src, p.Dest)
	return nil
}

func (b *DryRunBackend) FileMove(ctx context.Context, p *FileMoveParams) error {
	ui.Printf("mv %s %s\n", p.Src, p.Dest)
	return nil
}

func (b *DryRunBackend) Symlink(ctx context.Context, p *SymlinkParams) error {
	ui.Printf("ln -s %s %s\n", p.Src, p.Dest)
	return nil
}

func (b *DryRunBackend) HTTPRequest(ctx context.Context, p *HTTPRequestParams) error {
	ui.Printf("curl -Lo %s %s\n", p.Path, p.URL)
	return nil
}

func (b *DryRunBackend) Command(ctx context.Context, p *CommandParams) error {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s ", p.CmdName)
	for _, arg := range p.CmdArgs {
		fmt.Fprintf(buf, "%s ", arg)
	}
	fmt.Fprintf(buf, "\n")

	ui.Printf(buf.String())
	return nil
}
