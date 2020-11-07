//go:generate mockery -outpkg=mocks -case=snake -name=OSDetector

package modules

import (
	"bytes"
	"context"
	"os"
	"strings"
	"time"

	"github.com/wesovilabs/koazee"

	"github.com/raba-jp/primus/pkg/ctxlib"
	"go.uber.org/zap"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/spf13/afero"
)

const timeout = 5 * time.Second

var _ OSDetector = (*osDetector)(nil)

type OSDetector interface {
	ExecutableCommand(ctx context.Context, name string) (ok bool)
	Darwin(ctx context.Context) (v bool)
	ArchLinux(ctx context.Context) (v bool)
}

type osDetector struct {
	OSDetector
	exc exec.Interface
	fs  afero.Fs
}

func NewOSDetector(exc exec.Interface, fs afero.Fs) OSDetector {
	return &osDetector{
		exc: exc,
		fs:  fs,
	}
}

func (d *osDetector) Darwin(ctx context.Context) bool {
	ctx, _ = ctxlib.LoggerWithNamespace(ctx, "os_detector")
	ctx, logger := ctxlib.LoggerWithNamespace(ctx, "darwin")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	buferr := new(bytes.Buffer)
	cmd := d.exc.CommandContext(ctx, "uname", "-a")
	cmd.SetStderr(buferr)
	out, err := cmd.Output()
	if err != nil {
		logger.Error("Failed to detect darwin", zap.String("stderr", buferr.String()), zap.Error(err))
		return false
	}
	return strings.Contains(string(out), "Darwin")
}

func (d *osDetector) ArchLinux(ctx context.Context) bool {
	ctx, _ = ctxlib.LoggerWithNamespace(ctx, "os_detector")
	_, logger := ctxlib.LoggerWithNamespace(ctx, "arch_linux")
	_, err := d.fs.Stat("/etc/arch-release")
	if err != nil {
		logger.Debug("FS stats failed", zap.Error(err))
	}
	return err == nil
}

func (d *osDetector) ExecutableCommand(ctx context.Context, name string) bool {
	ctx, _ = ctxlib.LoggerWithNamespace(ctx, "os_detector")
	_, logger := ctxlib.LoggerWithNamespace(ctx, "executable_command")

	path := Getenv("PATH")
	paths := strings.Split(path, ":")
	logger.Debug("PATH environment variable", zap.Strings("paths", paths))

	executable := false
	koazee.StreamOf(paths).RemoveDuplicates().ForEach(func(path string) {
		_ = afero.Walk(d.fs, path, func(path string, _ os.FileInfo, _ error) error {
			executable = executable || strings.HasSuffix(path, name)
			return nil
		})
	})

	zap.L().Debug(
		"check executable",
		zap.String("name", name),
		zap.Strings("path", paths),
		zap.Bool("ok", executable),
	)

	return executable
}
