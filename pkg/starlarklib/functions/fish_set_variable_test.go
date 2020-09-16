package functions_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/internal/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/internal/handlers/mock"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"go.starlark.net/starlark"
)

func TestFishSetVariable(t *testing.T) {
	tests := []struct {
		name   string
		expr   string
		mock   func(*mock_handlers.MockFishSetVariableHandler)
		hasErr bool
	}{
		{
			name: "success",
			expr: `fish_set_variable(name="GOPATH", value="$HOME/go", scope="universal", export=True)`,
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
			name:   "error: too many arguments",
			expr:   `fish_set_variable("GOPATH", "$HOME/go", "universal", True, "too many")`,
			mock:   func(m *mock_handlers.MockFishSetVariableHandler) {},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockFishSetVariableHandler(ctrl)
			tt.mock(m)

			predeclared := starlark.StringDict{
				"fish_set_variable": starlark.NewBuiltin("fish_set_variable", functions.FishSetVariable(m)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.expr, predeclared)
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
