package os_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/functions/os"
	"github.com/raba-jp/primus/pkg/modules"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func TestNewIsDarwinFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mock modules.DarwinChecker
		want lib.Value
	}{
		{
			name: "success",
			mock: func(context.Context) bool {
				return true
			},
			want: lib.True,
		},
		{
			name: "failure",
			mock: func(context.Context) bool {
				return false
			},
			want: lib.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			detector := new(mocks.OSDetector)
			detector.ApplyDarwinExpectation(tt.mock)

			globals, err := starlark.ExecForTest("test", `v = test()`, os.NewIsDarwinFunction(detector))

			assert.NoError(t, err)
			assert.Equal(t, tt.want, globals["v"])
		})
	}
}

func TestNewDarwinInstalledFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      os.DarwinInstalledRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: true",
			data: `test(name="base-devel")`,
			mock: func(context.Context, string) bool {
				return true
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: false",
			data: `test(name="base-devel")`,
			mock: func(context.Context, string) bool {
				return false
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("base-devel", "too many")`,
			mock: func(context.Context, string) bool {
				return true
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, os.NewDarwinInstalledFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestNewDarwinInstallFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      os.DarwinInstallRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="base-devel", option="option", cask=True)`,
			mock: func(context.Context, *os.DarwinInstallParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("base-devel", "option", True, "too many")`,
			mock: func(context.Context, *os.DarwinInstallParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: package install failed",
			data: `test(name="base-devel")`,
			mock: func(context.Context, *os.DarwinInstallParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, os.NewDarwinInstallFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestNewDarwinUninstallFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      os.DarwinUninstallRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: func(context.Context, string) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("base-devel", "too many")`,
			mock: func(context.Context, string) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: func(context.Context, string) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, os.NewDarwinUninstallFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestDarwinInstalled(t *testing.T) {
	tests := []struct {
		name string
		mock command.ExecuteRunner
		fs   func() afero.Fs
		want bool
	}{
		{
			name: "success: $(brew --prefix)/Celler",
			mock: func(ctx context.Context, p *command.Params) error {
				p.Stdout.Write([]byte("/opt"))
				return nil
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/opt/Celler/cat", []byte{}, 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "success: /opt/homebrew-cask/Caskroom",
			mock: func(ctx context.Context, p *command.Params) error {
				p.Stdout.Write([]byte("/opt"))
				return nil
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/opt/homebrew-cask/Caskroom/cat", []byte{}, 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "success: /usr/local/Caskroom",
			mock: func(ctx context.Context, p *command.Params) error {
				p.Stdout.Write([]byte("/opt"))
				return nil
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/usr/local/Caskroom/cat", []byte{}, 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "error: not found",
			mock: func(ctx context.Context, p *command.Params) error {
				p.Stdout.Write([]byte("/opt"))
				return nil
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				return fs
			},
			want: false,
		},
		{
			name: "error: command failed",
			mock: func(ctx context.Context, p *command.Params) error {
				return xerrors.New("dummy")
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				return fs
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkInstall := os.DarwinInstalled(tt.mock, tt.fs())
			res := checkInstall(context.Background(), "cat")
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestDarwinInstall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		execute   func() command.ExecuteRunner
		fs        func() afero.Fs
		params    *os.DarwinInstallParams
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			execute: func() command.ExecuteRunner {
				return func(context.Context, *command.Params) error {
					return nil
				}
			},
			fs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &os.DarwinInstallParams{
				Name:   "pkg",
				Option: "options",
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: cask",
			execute: func() command.ExecuteRunner {
				return func(context.Context, *command.Params) error {
					return nil
				}
			},
			fs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			params: &os.DarwinInstallParams{
				Name:   "pkg",
				Option: "options",
				Cask:   true,
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: already installed",
			execute: func() command.ExecuteRunner {
				return func(context.Context, *command.Params) error {
					return nil
				}
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/usr/local/Caskroom/pkg", []byte{}, 0o777)
				return fs
			},
			params: &os.DarwinInstallParams{
				Name:   "pkg",
				Option: "options",
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: install package failed",
			execute: func() command.ExecuteRunner {
				return func(context.Context, *command.Params) error {
					return xerrors.New("dummy")
				}
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				return fs
			},
			params: &os.DarwinInstallParams{
				Name:   "pkg",
				Option: "options",
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			install := os.DarwinInstall(tt.execute(), tt.fs())
			err := install(context.Background(), tt.params)
			tt.errAssert(t, err)
		})
	}
}

func TestNewUninstall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		execute   func() command.ExecuteRunner
		fs        func() afero.Fs
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			execute: func() command.ExecuteRunner {
				return func(context.Context, *command.Params) error {
					return nil
				}
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/usr/local/Caskroom/pkg", []byte{}, 0o777)
				return fs
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: not installed",
			execute: func() command.ExecuteRunner {
				return func(context.Context, *command.Params) error {
					return nil
				}
			},
			fs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: uninstall failed",
			execute: func() command.ExecuteRunner {
				called := false
				return func(context.Context, *command.Params) error {
					if called {
						return xerrors.New("dummy")
					}
					called = true
					return nil
				}
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/usr/local/Caskroom/pkg", []byte{}, 0o777)
				return fs
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			uninstall := os.DarwinUninstall(tt.execute(), tt.fs())
			err := uninstall(context.Background(), "pkg")
			tt.errAssert(t, err)
		})
	}
}
