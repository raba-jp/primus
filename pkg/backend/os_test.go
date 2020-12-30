package backend_test

import (
	"context"
	"testing"

	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewDarwinChecker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mock backend.Execute
		want bool
	}{
		{
			name: "success",
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				p.Stdout.Write([]byte("Darwin myMac.local 15.3.0 Darwin Kernel Version 15.3.0: Thu Dec 10 18:40:58 PST 2015; root:xnu-3248.30.4~1/RELEASE_X86_64 x86_64"))
				return nil
			},
			want: true,
		},
		{
			name: "fail: linux",
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				p.Stdout.Write([]byte("Linux HostName 5.7.19-2-MANJARO #1 SMP PREEMPT Fri Aug 28 20:22:12 UTC 2020 x86_64 GNU/Linux"))
				return nil
			},
			want: false,
		},
		{
			name: "fail: command error",
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return xerrors.New("dummy")
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := backend.NewDarwinChecker(tt.mock)(context.Background())
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestNewArchLinuxChecker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mock func() afero.Fs
		want bool
	}{
		{
			name: "success",
			mock: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/etc/arch-release", []byte("Arch Linux"), 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "success: empty file",
			mock: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/etc/arch-release", []byte(""), 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "fail: not exists /etc/arch-release",
			mock: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := tt.mock()

			result := backend.NewArchLinuxChecker(fs)(context.Background())
			assert.Equal(t, tt.want, result)
		})
	}
}
