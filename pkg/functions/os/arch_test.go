package os_test

import (
	"context"
	"fmt"
	"testing"

	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/functions/os"

	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/stretchr/testify/assert"
	lib "go.starlark.net/starlark"
)

func TestNewIsArchFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mock backend.ArchLinuxChecker
		want lib.Value
	}{
		{
			name: "success",
			mock: func(ctx context.Context) bool {
				return true
			},
			want: lib.True,
		},
		{
			name: "success: returns false",
			mock: func(ctx context.Context) bool {
				return false
			},
			want: lib.False,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			globals, err := starlark.ExecForTest("test", `v = test()`, os.NewIsArchFunction(tt.mock))

			assert.NoError(t, err)
			assert.Equal(t, tt.want, globals["v"])
		})
	}
}

func TestNewArchInstalledFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      os.ArchInstalledRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: true",
			data: `test(name="base-devel")`,
			mock: func(ctx context.Context, name string) bool {
				return true
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: false",
			data: `test(name="base-devel")`,
			mock: func(ctx context.Context, name string) bool {
				return false
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("base-devel", "too many")`,
			mock: func(ctx context.Context, name string) bool {
				return true
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, os.NewArchInstalledFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestNewArchInstallFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      os.ArchInstallRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="base-devel", option="option")`,
			mock: func(ctx context.Context, p *os.ArchInstallParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("base-devel", "option", "too many")`,
			mock: func(ctx context.Context, p *os.ArchInstallParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: package install failed",
			data: `test(name="base-devel")`,
			mock: func(ctx context.Context, p *os.ArchInstallParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, os.NewArchInstallFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestNewArchMultipleInstallFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      os.ArchMultipleInstallRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(names=["arg1", "arg2"])`,
			mock: func(ctx context.Context, ps []*os.ArchInstallParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success, includes number",
			data: `test(names=["arg1", 1, "arg2"])`,
			mock: func(ctx context.Context, ps []*os.ArchInstallParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test(["arg1", "arg2"], "too many")`,
			mock: func(ctx context.Context, ps []*os.ArchInstallParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: package install failed",
			data: `test(names=["arg1", "arg2"])`,
			mock: func(ctx context.Context, ps []*os.ArchInstallParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, os.NewArchMultipleInstallFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestNewArchUninstallFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		data   string
		mock   os.ArchUninstallRunner
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: func(ctx context.Context, name string) error {
				return nil
			},
			hasErr: false,
		},
		{
			name: "error: too many arguments",
			data: `test("base-devel", "yay", "too many")`,
			mock: func(ctx context.Context, name string) error {
				return nil
			},
			hasErr: true,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: func(ctx context.Context, name string) error {
				return xerrors.New("dummy")
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, os.NewArchUninstallFunction(tt.mock))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestArchInstalled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		executable backend.Executable
		execute    backend.Execute
		want       bool
	}{
		{
			name: "success: returns true",
			executable: func(ctx context.Context, cmd string) bool {
				return true
			},
			execute: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			want: true,
		},
		{
			name: "success: returns false",
			executable: func(ctx context.Context, cmd string) bool {
				return true
			},
			execute: func(ctx context.Context, p *backend.ExecuteParams) error {
				return xerrors.New("dummy")
			},
			want: false,
		},
		{
			name: "success: cannot use yay",
			executable: func(ctx context.Context, cmd string) bool {
				return false
			},
			execute: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			checkInstall := os.ArchInstalled(tt.executable, tt.execute)
			res := checkInstall(context.Background(), "base-devel")
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestArchInstall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		executable backend.Executable
		execute    backend.Execute
		errAssert  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			executable: func(ctx context.Context, name string) bool {
				return false
			},
			execute: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: already installed",
			executable: func(ctx context.Context, name string) bool {
				return true
			},
			execute: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: install failed",
			executable: func(ctx context.Context, name string) bool {
				return false
			},
			execute: func(ctx context.Context, p *backend.ExecuteParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			install := os.ArchInstall(tt.executable, tt.execute)
			err := install(context.Background(), &os.ArchInstallParams{
				Name:   "base-devel",
				Option: "options",
			})
			tt.errAssert(t, err)
		})
	}
}

func TestNewMultipleInstall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		executable backend.Executable
		execute    backend.Execute
		errAssert  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			executable: func(ctx context.Context, name string) bool {
				return false
			},
			execute: func(context.Context, *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "failed",
			executable: func(context.Context, string) bool {
				return false
			},
			execute: func(context.Context, *backend.ExecuteParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			multipleInstall := os.ArchMultipleInstall(tt.executable, tt.execute)
			err := multipleInstall(context.Background(), []*os.ArchInstallParams{
				{Name: "arg1"},
				{Name: "arg2"},
			})
			tt.errAssert(t, err)
		})
	}
}

func TestArchUninstall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		executable backend.Executable
		mock       func() backend.Execute
		errAssert  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			executable: func(context.Context, string) bool {
				return true
			},
			mock: func() backend.Execute {
				return func(context.Context, *backend.ExecuteParams) error {
					return nil
				}
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: already installed",
			executable: func(context.Context, string) bool {
				return true
			},
			mock: func() backend.Execute {
				return func(context.Context, *backend.ExecuteParams) error {
					return xerrors.New("dummy")
				}
			},
			errAssert: assert.NoError,
		},
		{
			name: "failure",
			executable: func(context.Context, string) bool {
				return true
			},
			mock: func() backend.Execute {
				called := false
				return func(context.Context, *backend.ExecuteParams) error {
					fmt.Println(called)
					if called {
						return xerrors.New("dummy")
					}
					called = true
					return nil
				}
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uninstall := os.ArchUninstall(tt.executable, tt.mock())
			err := uninstall(context.Background(), "base-devel")
			tt.errAssert(t, err)
		})
	}
}

func TestArchCmdArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params []string
		mock   func() backend.Executable
		cmd    string
		args   []string
	}{
		{
			name:   "install, yay, powerpill",
			params: []string{"arg1", "arg2"},
			mock: func() backend.Executable {
				return func(context.Context, string) bool {
					return true
				}
			},
			cmd:  "yay",
			args: []string{"--pacman", "powerpill", "-S", "--noconfirm", "arg1", "arg2"},
		},
		{
			name:   "install, yay",
			params: []string{"arg1", "arg2"},
			mock: func() backend.Executable {
				called := false
				return func(context.Context, string) bool {
					if called {
						return false
					}
					called = true
					return true
				}
			},
			cmd:  "yay",
			args: []string{"-S", "--noconfirm", "arg1", "arg2"},
		},
		{
			name:   "install, powerpill",
			params: []string{"arg1", "arg2"},
			mock: func() backend.Executable {
				called := false
				return func(context.Context, string) bool {
					if called {
						return true
					}
					called = true
					return false
				}
			},
			cmd:  "sudo pacman",
			args: []string{"-S", "--noconfirm", "arg1", "arg2"},
		},
		{
			name:   "install, invalid args",
			params: []string{"arg1", "", "arg2"},
			mock: func() backend.Executable {
				return func(context.Context, string) bool {
					return false
				}
			},
			cmd:  "sudo pacman",
			args: []string{"-S", "--noconfirm", "arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd, args := os.ArchCmdArgs(context.Background(), tt.mock(), tt.params)
			assert.Equal(t, tt.cmd, cmd)
			assert.Equal(t, tt.args, args)
		})
	}
}
