package starlarkfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/command/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/command/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

func TestExecutable(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(m *mock_handlers.MockExecutableHandler)
		want   lib.Value
		hasErr bool
	}{
		{
			name: "success: return true",
			data: `v = test("data")`,
			mock: func(m *mock_handlers.MockExecutableHandler) {
				m.EXPECT().Executable(gomock.Any(), gomock.Any()).Return(true)
			},
			want:   lib.True,
			hasErr: false,
		},
		{
			name: "success: return false",
			data: `v = test("data")`,
			mock: func(m *mock_handlers.MockExecutableHandler) {
				m.EXPECT().Executable(gomock.Any(), gomock.Any()).Return(false)
			},
			want:   lib.False,
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `v = test("data", "too many")`,
			mock:   func(m *mock_handlers.MockExecutableHandler) {},
			want:   nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockExecutableHandler(ctrl)
			tt.mock(m)

			globals, err := starlark.ExecForTest("test", tt.data, starlarkfn.Executable(m))
			if !tt.hasErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if diff := cmp.Diff(globals["v"], tt.want); diff != "" {
				t.Error(diff)
			}
		})
	}
}
