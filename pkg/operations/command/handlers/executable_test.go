package handlers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

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
			name: "success: not found",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env.SetFakeEnv(tt.env)
			executable := handlers.NewExecutable(tt.fs())
			ret := executable.Run(context.Background(), tt.data)
			assert.Equal(t, tt.want, ret)
		})
	}
}
