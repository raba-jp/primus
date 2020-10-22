package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/command/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/command/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

func TestExecutable(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      func(m *mock_handlers.MockExecutableHandler)
		want      lib.Value
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: return true",
			data: `v = test("data")`,
			mock: func(m *mock_handlers.MockExecutableHandler) {
				m.EXPECT().Executable(gomock.Any(), gomock.Any()).Return(true)
			},
			want:      lib.True,
			errAssert: assert.NoError,
		},
		{
			name: "success: return false",
			data: `v = test("data")`,
			mock: func(m *mock_handlers.MockExecutableHandler) {
				m.EXPECT().Executable(gomock.Any(), gomock.Any()).Return(false)
			},
			want:      lib.False,
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `v = test("data", "too many")`,
			mock:      func(m *mock_handlers.MockExecutableHandler) {},
			want:      nil,
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockExecutableHandler(ctrl)
			tt.mock(m)

			globals, err := starlark.ExecForTest("test", tt.data, starlarkfn.Executable(m))
			tt.errAssert(t, err)
			assert.Equal(t, globals["v"], tt.want)
		})
	}
}
