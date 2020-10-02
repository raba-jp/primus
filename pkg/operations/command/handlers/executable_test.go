package handlers_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/env"

	"github.com/raba-jp/primus/pkg/operations/command/handlers"
	"github.com/spf13/afero"
)

func TestNewExecutable(t *testing.T) {
	tests := []struct {
		name string
		data string
		env  map[string]string
		fs   func() afero.Fs
		want bool
	}{
		{
			name: "success",
			data: "cat",
			env: map[string]string{
				"SHELL": "/bin/bash",
				"PATH":  "/bin:/usr/bin",
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/bin/cat", []byte{}, 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "success: fish shell",
			data: "cat",
			env: map[string]string{
				"SHELL": "/bin/local/bin/fish",
				"PATH":  "/bin /usr/bin /usr/local/bin",
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/bin/cat", []byte{}, 0o777)
				return fs
			},
			want: true,
		},
		{
			name: "success: posix shell with invalid PATH",
			data: "cat",
			env: map[string]string{
				"SHELL": "/bin/bash",
				"PATH":  "/bin /usr/bin /usr/local/bin",
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/bin/cat", []byte{}, 0o777)
				return fs
			},
			want: false,
		},
		{
			name: "success: fish shell with invalid PATH",
			data: "cat",
			env: map[string]string{
				"SHELL": "/bin/local/bin/fish",
				"PATH":  "/bin:/usr/bin:/usr/local/bin",
			},
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				afero.WriteFile(fs, "/bin/cat", []byte{}, 0o777)
				return fs
			},
			want: false,
		},
		{
			name: "success: posix shell, not found",
			data: "cat",
			env: map[string]string{
				"SHELL": "/bin/bash",
				"PATH":  "/bin:/usr/bin:/usr/local/bin",
			},
			fs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			want: false,
		},
		{
			name: "success: fish shell, not found",
			data: "cat",
			env: map[string]string{
				"SHELL": "/bin/local/bin/fish",
				"PATH":  "/bin /usr/bin /usr/local/bin",
			},
			fs: func() afero.Fs {
				return afero.NewMemMapFs()
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env.SetFakeEnv(tt.env)
			handler := handlers.NewExecutable(tt.fs())
			if ret := handler.Executable(context.Background(), tt.data); ret != tt.want {
				t.Errorf("Unexpected error: want: %v, got: %v", tt.want, ret)
			}
		})
	}
}
