package starlarkfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/operations/fish/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/fish/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/fish/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestSetVariable(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockSetVariableHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="GOPATH", value="$HOME/go", scope="universal", export=True)`,
			mock: func(m *mock_handlers.MockSetVariableHandler) {
				m.EXPECT().SetVariable(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.SetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.UniversalScope,
						Export: true,
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: args",
			data: `test("GOPATH", "$HOME/go", "universal", True)`,
			mock: func(m *mock_handlers.MockSetVariableHandler) {
				m.EXPECT().SetVariable(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.SetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.UniversalScope,
						Export: true,
					}),
				)
			},
			hasErr: false,
		},
		{
			name: "success: global scope",
			data: `test("GOPATH", "$HOME/go", "global", True)`,
			mock: func(m *mock_handlers.MockSetVariableHandler) {
				m.EXPECT().SetVariable(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.SetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.GlobalScope,
						Export: true,
					}),
				)
			},
			hasErr: false,
		},
		{
			name: "success: local scope",
			data: `test("GOPATH", "$HOME/go", "local", True)`,
			mock: func(m *mock_handlers.MockSetVariableHandler) {
				m.EXPECT().SetVariable(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.SetVariableParams{
						Name:   "GOPATH",
						Value:  "$HOME/go",
						Scope:  handlers.LocalScope,
						Export: true,
					}),
				)
			},
			hasErr: false,
		},
		{
			name:   "error: unexpected scope",
			data:   `test(name="GOPATH", value="$HOME/go", scope="dummy", export=True)`,
			mock:   func(m *mock_handlers.MockSetVariableHandler) {},
			hasErr: true,
		},
		{
			name:   "error: too many arguments",
			data:   `test("GOPATH", "$HOME/go", "universal", True, "too many")`,
			mock:   func(m *mock_handlers.MockSetVariableHandler) {},
			hasErr: true,
		},
		{
			name: "error: return handler error",
			data: `test("GOPATH", "$HOME/go", "universal", True)`,
			mock: func(m *mock_handlers.MockSetVariableHandler) {
				m.EXPECT().SetVariable(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockSetVariableHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.SetVariable(m))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
