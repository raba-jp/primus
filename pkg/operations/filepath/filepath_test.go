package filepath_test

import (
	"testing"

	lib "go.starlark.net/starlark"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/operations/filepath"
	"github.com/raba-jp/primus/pkg/starlark"
)

func TestCurrent(t *testing.T) {
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
				"test": lib.NewBuiltin("test", filepath.Current()),
			}
			globals, err := lib.ExecFile(starlark.NewThread("test"), tt.filepath, tt.data, predeclared)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			got, _ := lib.AsString(globals["v"])
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDir(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   string
		hasErr bool
	}{
		{
			name:   "success",
			data:   `v = test(path="/test.star")`,
			want:   "/",
			hasErr: false,
		},
		{
			name:   "success: has parent dir",
			data:   `v = test(path="/sym/test.star")`,
			want:   "/sym",
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `v = test("/sym/test.star", "too many")`,
			want:   "",
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globals, err := starlark.ExecForTest("test", tt.data, filepath.Dir())
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
			globals, err := starlark.ExecForTest("test", tt.data, filepath.Join())
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
