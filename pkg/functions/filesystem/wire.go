//+build wireinject

package filesystem

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/modules"
	lib "go.starlark.net/starlark"
)

func newFunctions(
	createDirectory CreateDirectoryRunner,
	createSymlink CreateSymlinkRunner,
	copyFile CopyFileRunner,
	moveFile MoveFileRunner,
	existsFile ExistsFileRunner,
) lib.Value {
	dict := lib.NewDict(5)
	dict.SetKey(lib.String("create_directory"), lib.NewBuiltin(
		"create_directory",
		NewCreateDirectoryFunction(createDirectory),
	))
	dict.SetKey(lib.String("create_symlink"), lib.NewBuiltin(
		"create_symlink",
		NewCreateSymlinkFunction(createSymlink),
	))
	dict.SetKey(lib.String("copy_file"), lib.NewBuiltin(
		"copy_file",
		NewCopyFileFunction(copyFile),
	))
	dict.SetKey(lib.String("move_file"), lib.NewBuiltin(
		"move_file",
		NewMoveFileFunction(moveFile),
	))
	dict.SetKey(lib.String("exists_file"), lib.NewBuiltin(
		"exists_file",
		NewExistsFileFunction(existsFile),
	))
	return dict
}

func NewFunctions() lib.Value {
	wire.Build(
		modules.NewFs,
		CreateDirectory,
		CreateSymlink,
		CopyFile,
		MoveFile,
		ExistsFile,
		newFunctions,
	)
	return nil
}
