package modules_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/modules"
	"github.com/stretchr/testify/assert"
)

func TestSetGlobalLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level     string
		errAssert assert.ErrorAssertionFunc
	}{
		{level: "debug", errAssert: assert.NoError},
		{level: "info", errAssert: assert.NoError},
		{level: "warn", errAssert: assert.NoError},
		{level: "error", errAssert: assert.NoError},
		{level: "test", errAssert: assert.Error},
		{level: "", errAssert: assert.Error},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			t.Parallel()
			modules.SetGlobalLogger(tt.level)
		})
	}
}
