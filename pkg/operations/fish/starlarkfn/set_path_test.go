package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/operations/fish/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/fish/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/fish/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestSetPath(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      func(*mock_handlers.MockSetPathHandler)
		errAssert assert.ErrorAssertionFunc
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
			errAssert: assert.NoError,
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
			errAssert: assert.NoError,
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
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test(["$GOPATH/bin", "$HOME/.bin"], "too many")`,
			mock:      func(m *mock_handlers.MockSetPathHandler) {},
			errAssert: assert.Error,
		},
		{
			name: "error: return handler error",
			data: `test(["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(m *mock_handlers.MockSetPathHandler) {
				m.EXPECT().SetPath(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockSetPathHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.SetPath(m))
			tt.errAssert(t, err)
		})
	}
}
