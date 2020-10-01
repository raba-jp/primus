//go:generate mockgen -destination mock/handler.go . InstallExtensionHandler

package handlers

import (
	"bytes"
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type InstallExtensionParams struct {
	Name    string
	Version string
}

type InstallExtensionHandler interface {
	InstallExtension(ctx context.Context, dryrun bool, p *InstallExtensionParams) error
}

type InstallExtensionHandlerFunc func(ctx context.Context, dryrun bool, p *InstallExtensionParams) error

func (f InstallExtensionHandlerFunc) InstallExtension(ctx context.Context, dryrun bool, p *InstallExtensionParams) error {
	return f(ctx, dryrun, p)
}

func NewInstallExtension(execIF exec.Interface) InstallExtensionHandler {
	return InstallExtensionHandlerFunc(func(ctx context.Context, dryrun bool, p *InstallExtensionParams) error {
		if dryrun {
			if p.Version != "" {
				// With version
				ui.Printf("code --install-extension %s@%s\n", p.Name, p.Version)
			} else {
				// Without version
				ui.Printf("code --install-extension %s\n", p.Name)
			}
			return nil
		}

		arg := p.Name
		if p.Version != "" {
			arg = arg + "@" + p.Version
		}

		cmd := execIF.CommandContext(ctx, "code", "--install-extension", arg)
		buf := new(bytes.Buffer)
		errbuf := new(bytes.Buffer)
		cmd.SetStdout(buf)
		cmd.SetStderr(errbuf)
		if err := cmd.Run(); err != nil {
			return xerrors.Errorf("failed install extension: %s: %w", arg, err)
		}
		zap.L().Info(
			"install vscode extension",
			zap.String("name", p.Name),
			zap.String("version", p.Version),
			zap.String("stdout", buf.String()),
			zap.String("stderr", errbuf.String()),
		)
		return nil
	})
}
