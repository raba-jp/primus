package backend_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "success",
			data: "echo",
			want: true,
		},
		{
			name: "success: full path",
			data: "/bin/echo",
			want: true,
		},
		{
			name: "success: not found",
			data: "aaaa",
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ret := backend.NewExecutable()(context.Background(), tt.data)
			assert.Equal(t, tt.want, ret)
		})
	}
}

func TestNewExecute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    *backend.ExecuteParams
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			params: &backend.ExecuteParams{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with user",
			params: &backend.ExecuteParams{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
				User: "root",
			},
			errAssert: assert.Error,
		},
		{
			name: "success: with cwd",
			params: &backend.ExecuteParams{
				Cmd:  "echo",
				Args: []string{"hello", "world"},
				Cwd:  "/",
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with stdin",
			params: &backend.ExecuteParams{
				Cmd:   "echo",
				Args:  []string{"hello", "world"},
				Stdin: new(bytes.Buffer),
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with stdout",
			params: &backend.ExecuteParams{
				Cmd:    "echo",
				Args:   []string{"hello", "world"},
				Stdout: new(bytes.Buffer),
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with stderr",
			params: &backend.ExecuteParams{
				Cmd:    "echo",
				Args:   []string{"hello", "world"},
				Stderr: new(bytes.Buffer),
			},
			errAssert: assert.NoError,
		},
		{
			name: "failure",
			params: &backend.ExecuteParams{
				Cmd:  "xxxxx",
				Args: []string{"hello", "world"},
				User: "root",
				Cwd:  "/",
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := backend.NewExecute()(context.Background(), tt.params)
			tt.errAssert(t, err)
		})
	}
}
