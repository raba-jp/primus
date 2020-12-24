package network_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/functions/network"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewGitCloneFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      network.GitCloneRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success:",
			data: `test(url="https://example.com", path="/sym", branch="main")`,
			mock: func(ctx context.Context, p *network.GitCloneParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: failed to git clone",
			data: `test("https://example.com", "/sym", "main")`,
			mock: func(ctx context.Context, p *network.GitCloneParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
		{
			name: "error: too many arguments",
			data: `test("https://example.com", "/sym", "main", "too many")`,
			mock: func(ctx context.Context, p *network.GitCloneParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, network.NewGitCloneFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestClone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    *network.GitCloneParams
		mock      command.ExecuteRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			params: &network.GitCloneParams{
				URL:    "https://github.com/raba-jp/dotfiles",
				Path:   "/tmp/dotfiles",
				Branch: "master",
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with cwd",
			params: &network.GitCloneParams{
				URL:    "https://github.com/raba-jp/dotfiles",
				Path:   "dotfiles",
				Branch: "master",
				Cwd:    "/tmp",
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "failure",
			params: &network.GitCloneParams{
				URL:    "https://github.com/raba-jp/dotfiles",
				Path:   "dotfiles",
				Branch: "master",
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := network.GitClone(tt.mock)(context.Background(), tt.params)
			tt.errAssert(t, err)
		})
	}
}
