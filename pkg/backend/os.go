package backend

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/afero"
)

const timeout = 5 * time.Second

type DarwinChecker func(ctx context.Context) bool

func NewDarwinChecker(execute Execute) DarwinChecker {
	return func(ctx context.Context) bool {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		bufout := new(bytes.Buffer)
		buferr := new(bytes.Buffer)
		if err := execute(ctx, &ExecuteParams{
			Cmd:    "uname",
			Args:   []string{"-a"},
			Stdout: bufout,
			Stderr: buferr,
		}); err != nil {
			log.Error().
				Str("stdout", bufout.String()).
				Str("stderr", buferr.String()).
				Err(err).
				Msg("failed to detect darwin")
			return false
		}
		return strings.Contains(bufout.String(), "Darwin")
	}
}

type ArchLinuxChecker func(ctx context.Context) bool

func NewArchLinuxChecker(fs afero.Fs) ArchLinuxChecker {
	return func(ctx context.Context) bool {
		_, err := fs.Stat("/etc/arch-release")
		if err != nil {
			log.Debug().Err(err).Msg("Filesystem stats failed")
		}
		return err == nil
	}
}
