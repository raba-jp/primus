//go:generate mockery -outpkg=mocks -case=snake -name=OSDetector

package modules

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/raba-jp/primus/pkg/exec"
	"github.com/spf13/afero"
)

const timeout = 5 * time.Second

type DarwinChecker func(ctx context.Context) bool

func NewDarwinChecker(exc exec.Interface) DarwinChecker {
	return func(ctx context.Context) bool {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		buferr := new(bytes.Buffer)
		cmd := exc.CommandContext(ctx, "uname", "-a")
		cmd.SetStderr(buferr)
		out, err := cmd.Output()
		if err != nil {
			log.Ctx(ctx).Error().
				Str("stderr", buferr.String()).
				Err(err).
				Msg("failed to detect darwin")
			return false
		}
		return strings.Contains(string(out), "Darwin")
	}
}

type ArchLinuxChecker func(ctx context.Context) bool

func NewArchLinuxChecker(fs afero.Fs) ArchLinuxChecker {
	return func(ctx context.Context) bool {
		_, err := fs.Stat("/etc/arch-release")
		if err != nil {
			log.Ctx(ctx).Debug().Err(err).Msg("Filesystem stats failed")
		}
		return err == nil
	}
}
