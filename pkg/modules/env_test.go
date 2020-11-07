package modules_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/modules"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	t.Parallel()

	osEnv := modules.Getenv("HOME")
	assert.NotEmpty(t, osEnv)

	modules.SetFakeEnv(map[string]string{
		"GOPATH": "/home/username/go",
	})
	fakeEnv := modules.Getenv("GOPATH")
	assert.NotEmpty(t, fakeEnv)
}
