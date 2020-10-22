package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/git/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/git/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestClone(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      func(m *mock_handlers.MockCloneHandler)
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success:",
			data: `test(url="https://example.com", path="/sym", branch="main")`,
			mock: func(m *mock_handlers.MockCloneHandler) {
				m.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: failed to git clone",
			data: `test("https://example.com", "/sym", "main")`,
			mock: func(m *mock_handlers.MockCloneHandler) {
				m.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			errAssert: assert.Error,
		},
		{
			name:      "error: too many arguments",
			data:      `test("https://example.com", "/sym", "main", "too many")`,
			mock:      func(m *mock_handlers.MockCloneHandler) {},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockCloneHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Clone(m))
			tt.errAssert(t, err)
		})
	}
}
