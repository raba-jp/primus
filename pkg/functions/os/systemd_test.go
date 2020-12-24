package os_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/functions/os"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestNewSystemdEnableFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      os.SystemdEnableRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="dummy.service")`,
			mock: func(ctx context.Context, name string) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("dummy.service", "too many")`,
			mock: func(ctx context.Context, name string) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to service enable",
			data: `test(name="dummy.service")`,
			mock: func(ctx context.Context, name string) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, os.NewSystemdEnableFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestNewSystemdStartFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      os.SystemdStartRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="dummy.service")`,
			mock: func(ctx context.Context, name string) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("dummy.service", "too many")`,
			mock: func(ctx context.Context, name string) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to service start",
			data: `test(name="dummy.service")`,
			mock: func(ctx context.Context, name string) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, os.NewSystemdStartFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestEnableService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mock      command.ExecuteRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: check cmd returns enabled",
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: enabled fail",
			mock: func(ctx context.Context, p *command.Params) error {
				if p.Args[0] == "is-enabled" {
					return nil
				}
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := os.SystemdEnable(tt.mock)(context.Background(), "dummy.service")
			tt.errAssert(t, err)
		})
	}
}

func TestSystemdStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mock      command.ExecuteRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: check cmd returns active",
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: enabled fail",
			mock: func(ctx context.Context, p *command.Params) error {
				if p.Args[0] == "is-active" {
					return nil
				}
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := os.SystemdStart(tt.mock)(context.Background(), "dummy.service")
			tt.errAssert(t, err)
		})
	}
}
