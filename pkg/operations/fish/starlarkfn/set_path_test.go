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

func TestSetPath(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockSetPathHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(values=["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(m *mock_handlers.MockSetPathHandler) {
				m.EXPECT().SetPath(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.SetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					}),
				).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: args",
			data: `test(["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(m *mock_handlers.MockSetPathHandler) {
				m.EXPECT().SetPath(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.SetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					}),
				)
			},
			hasErr: false,
		},
		{
			name: "success: include int and bool",
			data: `test(["$GOPATH/bin", 1, True, "$HOME/.bin"])`,
			mock: func(m *mock_handlers.MockSetPathHandler) {
				m.EXPECT().SetPath(
					gomock.Any(),
					gomock.Any(),
					gomock.Eq(&handlers.SetPathParams{
						Values: []string{"$GOPATH/bin", "$HOME/.bin"},
					}),
				)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test(["$GOPATH/bin", "$HOME/.bin"], "too many")`,
			mock:   func(m *mock_handlers.MockSetPathHandler) {},
			hasErr: true,
		},
		{
			name: "error: return handler error",
			data: `test(["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(m *mock_handlers.MockSetPathHandler) {
				m.EXPECT().SetPath(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockSetPathHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.SetPath(m))
			if !tt.hasErr && err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}
