package handlers_test

import (
	"bytes"

	"github.com/raba-jp/primus/pkg/exec"
	fakeexec "github.com/raba-jp/primus/pkg/exec/testing"
)

func newFakeOutputScript(action fakeexec.FakeAction) fakeexec.FakeCommandAction {
	return fakeexec.FakeCommandAction(func(cmd string, args ...string) exec.Cmd {
		fake := &fakeexec.FakeCmd{
			Stdout:       new(bytes.Buffer),
			Stderr:       new(bytes.Buffer),
			OutputScript: []fakeexec.FakeAction{action},
		}
		return fakeexec.InitFakeCmd(fake, cmd, args...)
	})
}

func newFakeRunScript(action fakeexec.FakeAction) fakeexec.FakeCommandAction {
	return fakeexec.FakeCommandAction(func(cmd string, args ...string) exec.Cmd {
		fake := &fakeexec.FakeCmd{
			Stdout:    new(bytes.Buffer),
			Stderr:    new(bytes.Buffer),
			RunScript: []fakeexec.FakeAction{action},
		}
		return fakeexec.InitFakeCmd(fake, cmd, args...)
	})
}
