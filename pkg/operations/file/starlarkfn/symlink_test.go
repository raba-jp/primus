package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/file/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/file/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestSymlink(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      func(*mock_handlers.MockSymlinkHandler)
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(src="/sym/src.txt", dest="/sym/dest.txt")`,
			mock: func(m *mock_handlers.MockSymlinkHandler) {
				m.EXPECT().Symlink(
					gomock.Any(),
					gomock.Any(),
					&handlers.SymlinkParams{
						Src:  "/sym/src.txt",
						Dest: "/sym/dest.txt",
					},
				).Return(nil)
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("/sym/src.txt", "/sys/dest.txt", "too many")`,
			mock:      func(m *mock_handlers.MockSymlinkHandler) {},
			errAssert: assert.Error,
		},
		{
			name: "error: create symlink failed ",
			data: `test("/sym/src.txt", "/sys/dest.txt")`,
			mock: func(m *mock_handlers.MockSymlinkHandler) {
				m.EXPECT().Symlink(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockSymlinkHandler(ctrl)
			tt.mock(m)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Symlink(m))
			tt.errAssert(t, err)
		})
	}
}
