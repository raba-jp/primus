package os_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	lib "go.starlark.net/starlark"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/functions/os"
	"github.com/raba-jp/primus/pkg/starlark"
)

func TestCurrentPath(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		data     string
		want     string
	}{
		{
			name:     "success",
			filepath: "/test.star",
			data:     "v = test()",
			want:     "/test.star",
		},
		{
			name:     "success: has parent dir",
			filepath: "/sym/test.star",
			data:     "v = test()",
			want:     "/sym/test.star",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predeclared := lib.StringDict{
				"test": lib.NewBuiltin("test", os.GetCurrentPath()),
			}
			globals, err := lib.ExecFile(starlark.NewThread("test"), tt.filepath, tt.data, predeclared)
			assert.NoError(t, err)

			got, _ := lib.AsString(globals["v"])
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDir(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		want      string
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name:      "success",
			data:      `v = test(path="/test.star")`,
			want:      "/",
			errAssert: assert.NoError,
		},
		{
			name:      "success: has parent dir",
			data:      `v = test(path="/sym/test.star")`,
			want:      "/sym",
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `v = test("/sym/test.star", "too many")`,
			want:      "",
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globals, err := starlark.ExecForTest("test", tt.data, os.GetDir())
			tt.errAssert(t, err)

			got, _ := lib.AsString(globals["v"])
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   string
		hasErr bool
	}{
		{
			name:   "success",
			data:   `v = test(paths=["/sym", "test", "test.star"])`,
			want:   "/sym/test/test.star",
			hasErr: false,
		},
		{
			name:   "success: includes int value",
			data:   `v = test(paths=["/sym", "test", 1, "test.star"])`,
			want:   "/sym/test/test.star",
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `v = test(["sym", "test.star"], "too many")`,
			want:   "",
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globals, err := starlark.ExecForTest("test", tt.data, os.JoinPath())
			if !tt.hasErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			got, _ := lib.AsString(globals["v"])
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Error(diff)
			}
		})
	}
}
