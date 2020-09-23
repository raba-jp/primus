package backend

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const timeout = 10 * time.Minute

var _ Backend = (*BaseBackend)(nil)

type BaseBackend struct {
	Exec        exec.Interface
	Fs          afero.Fs
	Client      *http.Client
	FnCallCount int
}

func (b *BaseBackend) FunctionCallCount() int {
	return b.FnCallCount
}

func (b *BaseBackend) CheckInstall(ctx context.Context, name string) bool {
	panic("Delegate to other backend")
}

func (b *BaseBackend) Install(ctx context.Context, dryrun bool, p *handlers.InstallParams) error {
	panic("Delegate to other backend")
}

func (b *BaseBackend) Uninstall(ctx context.Context, dryrun bool, p *handlers.UninstallParams) error {
	panic("Delegate to other backend")
}

func (b *BaseBackend) FileCopy(ctx context.Context, dryrun bool, p *handlers.FileCopyParams) error {
	if dryrun {
		ui.Printf("cp %s %s\n", p.Src, p.Dest)
		return nil
	}

	srcFile, err := b.Fs.Open(p.Src)
	if err != nil {
		return xerrors.Errorf("Failed to open src file: %w", err)
	}
	destFile, err := b.Fs.OpenFile(p.Dest, os.O_WRONLY|os.O_CREATE, p.Permission)
	if err != nil {
		return xerrors.Errorf("Failed to open dest file: %w", err)
	}
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return xerrors.Errorf("Failed to copy src to dest: %w", err)
	}
	zap.L().Info(
		"copied file",
		zap.String("source", p.Src),
		zap.String("destination", p.Dest),
		zap.String("permission", fmt.Sprintf("%v", p.Permission)),
	)
	return nil
}

func (b *BaseBackend) FileMove(ctx context.Context, dryrun bool, p *handlers.FileMoveParams) error {
	if dryrun {
		ui.Printf("mv %s %s\n", p.Src, p.Dest)
		return nil
	}

	if err := b.Fs.Rename(p.Src, p.Dest); err != nil {
		return xerrors.Errorf("Failed to move file: %s => %s: %w", p.Src, p.Dest, err)
	}
	zap.L().Info(
		"moved file",
		zap.String("source", p.Src),
		zap.String("destination", p.Dest),
	)
	return nil
}

func (b *BaseBackend) Symlink(ctx context.Context, dryrun bool, p *handlers.SymlinkParams) error {
	if dryrun {
		ui.Printf("ln -s %s %s\n", p.Src, p.Dest)
		return nil
	}

	if ext := b.fileExists(p.Dest); ext {
		return xerrors.New("File already exists")
	}

	linker, ok := b.Fs.(afero.Symlinker)
	if !ok {
		return xerrors.New("This filesystem does not support symlink")
	}
	if err := linker.SymlinkIfPossible(p.Src, p.Dest); err != nil {
		return xerrors.Errorf("Failed to create symbolic link: %w", err)
	}

	zap.L().Info(
		"create symbolic link",
		zap.String("source", p.Src),
		zap.String("destination", p.Dest),
	)

	return nil
}

func (b *BaseBackend) HTTPRequest(ctx context.Context, dryrun bool, p *handlers.HTTPRequestParams) error {
	if dryrun {
		ui.Printf("curl -Lo %s %s\n", p.Path, p.URL)
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, p.URL, nil)
	if err != nil {
		return xerrors.Errorf("Failed to create new http request: %w", err)
	}
	req = req.WithContext(ctx)

	res, err := b.Client.Do(req)
	if err != nil {
		return xerrors.Errorf("Failed to http request: %w", err)
	}
	defer res.Body.Close()

	if err := afero.WriteReader(b.Fs, p.Path, res.Body); err != nil {
		return xerrors.Errorf("Failed to write response body: %w", err)
	}

	return nil
}

func (b *BaseBackend) fileExists(path string) bool {
	_, err := b.Fs.Stat(path)
	if err == nil {
		zap.L().Info("Already exists file")
		return true
	}
	return false
}

func (b *BaseBackend) FishSetVariable(ctx context.Context, dryrun bool, p *handlers.FishSetVariableParams) error {
	var scope string
	switch p.Scope {
	case handlers.FishVariableUniversalScope:
		scope = "--universal"
	case handlers.FishVariableGlobalScope:
		scope = "--global"
	case handlers.FishVariableLocalScope:
		scope = "--local"
	}

	export := ""
	if p.Export {
		export = " --export"
	}

	arg := fmt.Sprintf("'set %s%s %s %s'", scope, export, p.Name, p.Value)

	if dryrun {
		ui.Printf("fish --command %s\n", arg)
		return nil
	}

	cmd := b.Exec.CommandContext(ctx, "fish", "--command", arg)
	buf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	cmd.SetStdout(buf)
	cmd.SetStderr(errbuf)
	if err := cmd.Run(); err != nil {
		return xerrors.Errorf("failed to set variable: fish --command %s: %w", arg, err)
	}
	zap.L().Info(
		"set fish variable",
		zap.String("name", p.Name),
		zap.String("value", p.Value),
		zap.String("scope", scope),
		zap.Bool("export", p.Export),
		zap.String("stdout", buf.String()),
		zap.String("stderr", errbuf.String()),
	)
	return nil
}

func (b *BaseBackend) FishSetPath(ctx context.Context, dryrun bool, p *handlers.FishSetPathParams) error {
	path := fmt.Sprintf("'set --universal fish_user_paths %s'", strings.Join(p.Values, " "))

	if dryrun {
		ui.Printf("fish --command %s\n", path)
		return nil
	}

	cmd := b.Exec.CommandContext(ctx, "fish", "--command", path)
	buf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	cmd.SetStdout(buf)
	cmd.SetStderr(errbuf)
	if err := cmd.Run(); err != nil {
		return xerrors.Errorf("failed to set path: fish --command 'set --universal fish_user_path %s': %w", path, err)
	}
	zap.L().Info("set fish path", zap.Strings("values", p.Values))
	return nil
}

func (b *BaseBackend) CreateDirectory(ctx context.Context, dryrun bool, p *handlers.CreateDirectoryParams) error {
	if dryrun {
		ui.Printf("mkdir -p %s\n", p.Path)
		ui.Printf("chmod %o %s\n", p.Permission, p.Path)
		return nil
	}

	if err := b.Fs.MkdirAll(p.Path, p.Permission); err != nil {
		return xerrors.Errorf("Create directory fialed: %w", err)
	}
	zap.L().Info("create directory", zap.String("path", p.Path), zap.String("permission", p.Permission.String()))
	return nil
}
