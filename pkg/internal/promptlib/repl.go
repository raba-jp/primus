package promptlib

import (
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/starlarklib/repl"
)

// const (
// 	defaultPrefix  = ">>> "
// 	enteringPrefix = "... "
// )
//
// var entring = false

func NewPrompt() {
	p := prompt.New(
		Executor(),
		NewCompeter,
	//prompt.OptionPrefix(defaultPrefix),
	//prompt.OptionLivePrefix(func() (string, bool) {
	//	if entring {
	//		return enteringPrefix, true
	//	}
	//	return defaultPrefix, false
	//})
	)
	p.Run()
}

func NewCompeter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "execute", Description: "Execute command. `execute('cmd', ['option1', 'option2'])`"},
		{Text: "file_copy", Description: "File copy. `file_copy('/a/src.txt', '/b/dest.txt')`"},
		{Text: "file_move", Description: "File move. `file_move('/a/src/txt', '/b/dest.txt')`"},
		{Text: "fish_set_path", Description: "Set $PATH. `fish_set_path(['$GOPATH/bin', '$HOME/.bin'])`"},
		{Text: "fish_set_varialbe", Description: "Set variable. `fish_set_variable('GOPATH', '$HOME/go')`"},
		{Text: "http_request", Description: "Send HTTP request. `http_request('https://example.com', '$HOME/example.html')`"},
		{Text: "package", Description: "Install package. `package('base-devel')`"},
		{Text: "symlink", Description: "Create symbolic link. `symlink('/a/src.txt', '/b/dest.txt')`"},
	}

	return prompt.FilterFuzzy(s, d.GetWordBeforeCursor(), true)
}

func Executor() func(s string) {
	repl := repl.Initialize()
	return func(s string) {
		s = strings.TrimSpace(s)
		if s == "exit" || s == "quit" {
			os.Exit(0)
			return
		}

		if s == "" {
			return
		}

		if err := repl.Eval(s); err != nil {
			ui.Errorf("%v\n", err)
		}
	}
}