package starlarkfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_handlers "github.com/raba-jp/primus/pkg/operations/packages/handlers/mock"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestInstall(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   func(*mock_handlers.MockCheckInstallHandler, *mock_handlers.MockInstallHandler)
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="base-devel", option="option")`,
			mock: func(ch *mock_handlers.MockCheckInstallHandler, i *mock_handlers.MockInstallHandler) {
				ch.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(false)
				i.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			hasErr: false,
		},
		{
			name: "success: already installed",
			data: `test(name="base-devel", option="option")`,
			mock: func(ch *mock_handlers.MockCheckInstallHandler, i *mock_handlers.MockInstallHandler) {
				ch.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(true)
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "option", "too many")`,
			mock:   func(ch *mock_handlers.MockCheckInstallHandler, i *mock_handlers.MockInstallHandler) {},
			hasErr: true,
		},
		{
			name: "error: package install failed",
			data: `test(name="base-devel")`,
			mock: func(ch *mock_handlers.MockCheckInstallHandler, i *mock_handlers.MockInstallHandler) {
				ch.EXPECT().CheckInstall(gomock.Any(), gomock.Any()).Return(false)
				i.EXPECT().Install(gomock.Any(), gomock.Any(), gomock.Any()).Return(xerrors.New("dummy"))
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ch := mock_handlers.NewMockCheckInstallHandler(ctrl)
			i := mock_handlers.NewMockInstallHandler(ctrl)

			tt.mock(ch, i)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Install(ch, i))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
