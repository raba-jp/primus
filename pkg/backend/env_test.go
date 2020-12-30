package backend_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	t.Parallel()

	osEnv := backend.Getenv("HOME")
	assert.NotEmpty(t, osEnv)

	backend.SetFakeEnv(map[string]string{
		"GOPATH": "/home/username/go",
	})
	fakeEnv := backend.Getenv("GOPATH")
	assert.NotEmpty(t, fakeEnv)
}
