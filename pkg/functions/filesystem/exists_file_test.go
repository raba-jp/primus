package filesystem_test

import (
	"context"
	"testing"

	lib "go.starlark.net/starlark"

	"github.com/raba-jp/primus/pkg/starlark"

	"github.com/raba-jp/primus/pkg/functions/filesystem"
	"github.com/stretchr/testify/assert"
)

func TestNewExistsFileFunction(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      filesystem.ExistsFileRunner
		want      lib.Bool
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: returns true",
			data: `v = test(path="/sym")`,
			mock: func(ctx context.Context, path string) bool {
				return true
			},
			want:      lib.True,
			errAssert: assert.NoError,
		},
		{
			name: "success: returns false",
			data: `v = test(path="/sym")`,
			mock: func(ctx context.Context, path string) bool {
				return false
			},
			want:      lib.False,
			errAssert: assert.NoError,
		},
		{
			name: "failure: too many argument",
			data: `v = test("/sym", "too many")`,
			mock: func(ctx context.Context, path string) bool {
				return false
			},
			want:      lib.False,
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globals, err := starlark.ExecForTest("test", tt.data, filesystem.NewExistsFileFunction(tt.mock))
			tt.errAssert(t, err)
			if globals["v"] != nil {
				assert.Equal(t, globals["v"], tt.want)
			}
		})
	}
}
