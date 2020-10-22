package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"
	"github.com/raba-jp/primus/pkg/operations/systemd/handlers"
	"golang.org/x/xerrors"
)

func TestNewEnableService(t *testing.T) {
	successAction := fakeexec.FakeAction(func() ([]byte, []byte, error) {
		return []byte{}, []byte{}, nil
	})
	enabledAction := fakeexec.FakeAction(func() ([]byte, []byte, error) {
		return []byte("enabled"), []byte{}, nil
	})
	failureAction := fakeexec.FakeAction(func() ([]byte, []byte, error) {
		return []byte{}, []byte{}, xerrors.New("dummy")
	})

	tests := []struct {
		name      string
		mock      exec.Interface
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			mock: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					newFakeOutputScript(successAction),
					newFakeRunScript(successAction),
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: check cmd returns enabled",
			mock: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					newFakeOutputScript(enabledAction),
					newFakeRunScript(successAction),
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: check fail",
			mock: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					newFakeOutputScript(failureAction),
				},
			},
			errAssert: assert.Error,
		},
		{
			name: "error: enabled fail",
			mock: &fakeexec.FakeExec{
				CommandScript: []fakeexec.FakeCommandAction{
					newFakeOutputScript(successAction),
					newFakeRunScript(failureAction),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.NewEnableService(tt.mock)
			err := handler.EnableService(context.Background(), false, "dummy.service")
			tt.errAssert(t, err)
		})
	}
}

func TestNewEnableService__DryRun(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "success",
			in:   "dummy.service",
			out:  "systemctl enable dummy.service\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			ui.SetDefaultUI(&ui.CommandLine{Out: buf, Errout: buf})
			handler := handlers.NewEnableService(nil)
			if err := handler.EnableService(context.Background(), true, tt.in); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.out, buf.String()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
