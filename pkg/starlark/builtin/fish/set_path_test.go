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

func TestSetPath(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockFishSetPathHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(values=["$GOPATH/bin", "$HOME/.bin"])`,
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
			name: "success: args",
			data: `test(["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(m *mock_handlers.MockFishSetPathHandler) {
				m.EXPECT().FishSetPath(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FishSetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					}),
				)
			},
			hasErr: false,
		},
		{
			name: "success: include int and bool",
			data: `test(["$GOPATH/bin", 1, True, "$HOME/.bin"])`,
			mock: func(m *mock_handlers.MockFishSetPathHandler) {
				m.EXPECT().FishSetPath(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.FishSetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					}),
				)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test(["$GOPATH/bin", "$HOME/.bin"], "too many")`,
			mock:   func(m *mock_handlers.MockFishSetPathHandler) {},
			hasErr: true,
		},
		{
			name: "error: return handler error",
			data: `test(["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(m *mock_handlers.MockFishSetPathHandler) {
				m.EXPECT().FishSetPath(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockFishSetPathHandler(ctrl)
			tt.mock(m)

			_, err := builtin.ExecForTest("test", tt.data, fish.SetPath(m))
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
