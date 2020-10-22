package starlark_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	lib "go.starlark.net/starlark"
)

func TestLoad(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, "/sym/child.star", []byte(
		`
def child():
	return None
`), 0644)
	assert.NoError(t, err)

	thread := &lib.Thread{
		Name: "main",
		Load: starlark.Load(fs, nil),
	}

	_, err = lib.ExecFile(thread, "/sym/parent.star", `load("child.star", "child")`, nil)
	assert.NoError(t, err)
}

func TestLoad_Nested(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, "/sym/child.star", []byte(
		`
load("child2.star", "child2")
def child():
	child2()
	return None
`), 0644)
	assert.NoError(t, err)
	err = afero.WriteFile(fs, "/sym/child2.star", []byte(
		`
def child2():
	return None
`), 0644)
	assert.NoError(t, err)

	thread := &lib.Thread{
		Name: "main",
		Load: starlark.Load(fs, nil),
	}

	_, err = lib.ExecFile(thread, "/sym/parent.star", `
load("child.star", "child")
child()
`, nil)
	assert.NoError(t, err)
}

func TestLoad_Nested__FileNotExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, "/sym/child.star", []byte(
		`
load("child2.star", "child2")
def child():
	child2()
	return None
`), 0644)
	assert.NoError(t, err)

	thread := &lib.Thread{
		Name: "main",
		Load: starlark.Load(fs, nil),
	}

	_, err = lib.ExecFile(thread, "/sym/parent.star", `
load("child.star", "child")
child()
`, nil)
	assert.Error(t, err)
}

func TestLoad_AbstractPath(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, "/sym/child.star", []byte(
		`
def child():
	return None
`), 0644)
	assert.NoError(t, err)

	thread := &lib.Thread{
		Name: "main",
		Load: starlark.Load(fs, nil),
	}

	_, err = lib.ExecFile(thread, "/sym/parent.star", `load("/sym/child.star", "child")`, nil)
	assert.NoError(t, err)
}
