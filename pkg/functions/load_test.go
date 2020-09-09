package functions_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/functions"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
)

func TestLoad(t *testing.T) {
	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, "/sym/child.star", []byte(
		`
def child():
	return None
`), 0644); err != nil {
		t.Fatalf("%v", err)
	}

	thread := &starlark.Thread{
		Name: "main",
		Load: functions.Load(fs),
	}
	starlarklib.SetCtx(context.Background(), thread)

	_, err := starlark.ExecFile(thread, "/sym/parent.star", `load("child.star", "child")`, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestLoad_Nested(t *testing.T) {
	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, "/sym/child.star", []byte(
		`
load("child2.star", "child2")
def child():
	child2()
	return None
`), 0644); err != nil {
		t.Fatalf("%v", err)
	}
	if err := afero.WriteFile(fs, "/sym/child2.star", []byte(
		`
def child2():
	return None
`), 0644); err != nil {
		t.Fatalf("%v", err)
	}

	thread := &starlark.Thread{
		Name: "main",
		Load: functions.Load(fs),
	}
	starlarklib.SetCtx(context.Background(), thread)

	_, err := starlark.ExecFile(thread, "/sym/parent.star", `
load("child.star", "child")
child()
`, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestLoad_AbstractPath(t *testing.T) {
	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, "/sym/child.star", []byte(
		`
def child():
	return None
`), 0644); err != nil {
		t.Fatalf("%v", err)
	}

	thread := &starlark.Thread{
		Name: "main",
		Load: functions.Load(fs),
	}
	starlarklib.SetCtx(context.Background(), thread)

	_, err := starlark.ExecFile(thread, "/sym/parent.star", `load("/sym/child.star", "child")`, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
