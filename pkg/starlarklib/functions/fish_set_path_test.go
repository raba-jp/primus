package functions_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/handlers/mock"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"go.starlark.net/starlark"
)

func TestFishSetPath(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockFishSetPathHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `fish_set_path(values=["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(m *mock_handlers.MockFishSetPathHandler) {
				m.EXPECT().FishSetPath(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FishSetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `fish_set_path(["$GOPATH/bin", "$HOME/.bin"], "too many")`,
			mock:   func(m *mock_handlers.MockFishSetPathHandler) {},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockFishSetPathHandler(ctrl)
			tt.mock(m)

			predeclared := starlark.StringDict{
				"fish_set_path": starlark.NewBuiltin("fish_set_path", functions.FishSetPath(m)),
			}
			thread := &starlark.Thread{
				Name: "testing",
			}
			_, err := starlark.ExecFile(thread, "test.star", tt.data, predeclared)
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
