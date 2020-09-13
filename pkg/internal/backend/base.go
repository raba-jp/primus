package backend

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const timeout = 10 * time.Minute

var _ Backend = (*BaseBackend)(nil)

type BaseBackend struct {
	Exec   exec.Interface
	Fs     afero.Fs
	Client *http.Client
}

func (b *BaseBackend) CheckInstall(ctx context.Context, name string) bool {
	panic("Delegate to other backend")
}

func (b *BaseBackend) Install(ctx context.Context, p *InstallParams) error {
	panic("Delegate to other backend")
}

func (b *BaseBackend) Uninstall(ctx context.Context, name string) error {
	panic("Delegate to other backend")
}

func (b *BaseBackend) Command(ctx context.Context, p *CommandParams) error {
	cmd := b.Exec.CommandContext(ctx, p.CmdName, p.CmdArgs...)
	buf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	cmd.SetStdout(buf)
	cmd.SetStderr(errbuf)
	if p.Cwd != "" {
		zap.L().Debug("Set directory", zap.String("cwd", p.Cwd))
		cmd.SetDir(p.Cwd)
	}

	if p.User != "" {
		proc, err := newSysProcAttr(p.User)
		if err != nil {
			return err
		}
		cmd.SetSysProcAttr(proc)
	}

	if err := cmd.Run(); err != nil {
		return xerrors.Errorf(
			"Failed to execute command '%s %s': %w",
			p.CmdName,
			p.CmdArgs,
			err,
		)
	}
	zap.L().Debug(
		"command output",
		zap.String("stdout", buf.String()),
		zap.String("stderr", errbuf.String()),
	)
	zap.L().Info(
		"Executed command",
		zap.String("cmd", p.CmdName),
		zap.Strings("args", p.CmdArgs),
		zap.String("user", p.User),
		zap.String("cwd", p.Cwd),
	)
	return nil
}

func getUser(name string) (*user.User, error) {
	u, err := user.Lookup(name)
	if err != nil {
		return nil, xerrors.Errorf("Failed to lookup user: %w", err)
	}
	return u, nil
}

func getUID(u *user.User) (uint32, error) {
	uid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return 0, xerrors.Errorf("%w", err)
	}
	return uint32(uid), nil
}

func getGID(u *user.User) (uint32, error) {
	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return 0, xerrors.Errorf("%w", err)
	}
	return uint32(gid), nil
}

func newSysProcAttr(name string) (*syscall.SysProcAttr, error) {
	u, err := getUser(name)
	if err != nil {
		return nil, err
	}
	uid, err := getUID(u)
	if err != nil {
		return nil, err
	}
	gid, err := getGID(u)
	if err != nil {
		return nil, err
	}
	proc := &syscall.SysProcAttr{}
	proc.Credential = &syscall.Credential{Uid: uid, Gid: gid}

	return proc, nil
}

func (b *BaseBackend) FileCopy(ctx context.Context, p *FileCopyParams) error {
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

func (b *BaseBackend) FileMove(ctx context.Context, p *FileMoveParams) error {
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

func (b *BaseBackend) Symlink(ctx context.Context, p *SymlinkParams) error {
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

func (b *BaseBackend) HTTPRequest(ctx context.Context, p *HTTPRequestParams) error {
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
