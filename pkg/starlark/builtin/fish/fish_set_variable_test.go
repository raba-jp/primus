package fish_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/handlers/mock"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
	"github.com/raba-jp/primus/pkg/starlark/builtin/fish"
	"golang.org/x/xerrors"
)

func TestSetVariable(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockFishSetVariableHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="GOPATH", value="$HOME/go", scope="universal", export=True)`,
			mock: func(m *mock_handlers.MockFishSetVariableHandler) {
				m.EXPECT().FishSetVariable(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FishSetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.FishVariableUniversalScope,
						Export: true,
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: args",
			data: `test("GOPATH", "$HOME/go", "universal", True)`,
			mock: func(m *mock_handlers.MockFishSetVariableHandler) {
				m.EXPECT().FishSetVariable(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FishSetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.FishVariableUniversalScope,
						Export: true,
					}),
				)
			},
			hasErr: false,
		},
		{
			name: "success: global scope",
			data: `test("GOPATH", "$HOME/go", "global", True)`,
			mock: func(m *mock_handlers.MockFishSetVariableHandler) {
				m.EXPECT().FishSetVariable(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FishSetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.FishVariableGlobalScope,
						Export: true,
					}),
				)
			},
			hasErr: false,
		},
		{
			name: "success: local scope",
			data: `test("GOPATH", "$HOME/go", "local", True)`,
			mock: func(m *mock_handlers.MockFishSetVariableHandler) {
				m.EXPECT().FishSetVariable(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FishSetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.FishVariableLocalScope,
						Export: true,
					}),
				)
			},
			hasErr: false,
		},
		{
			name:   "error: unexpected scope",
			data:   `test(name="GOPATH", value="$HOME/go", scope="dummy", export=True)`,
			mock:   func(m *mock_handlers.MockFishSetVariableHandler) {},
			hasErr: true,
		},
		{
			name:   "error: too many arguments",
			data:   `test("GOPATH", "$HOME/go", "universal", True, "too many")`,
			mock:   func(m *mock_handlers.MockFishSetVariableHandler) {},
			hasErr: true,
		},
		{
			name: "error: return handler error",
			data: `test("GOPATH", "$HOME/go", "universal", True)`,
			mock: func(m *mock_handlers.MockFishSetVariableHandler) {
				m.EXPECT().FishSetVariable(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockFishSetVariableHandler(ctrl)
			tt.mock(m)

			err := builtin.ExecForTest("test", tt.data, fish.SetVariable(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
