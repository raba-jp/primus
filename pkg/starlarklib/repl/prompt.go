package repl

import (
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/raba-jp/primus/pkg/cli/ui"
)

const (
	defaultPrefix  = ">>> "
	enteringPrefix = "... "
)

type PromptFunc func()

func NewPrompt(repl REPL, executor prompt.Executor, completer prompt.Completer) PromptFunc {
	return func() {
		p := prompt.New(
			executor,
			completer,
			prompt.OptionPrefix(defaultPrefix),
			prompt.OptionLivePrefix(func() (string, bool) {
				if repl.IsContinuation() {
					return enteringPrefix, true
				}
				return defaultPrefix, false
			}),
		)
		p.Run()
	}
}

func NewCompleter() prompt.Completer {
	return func(d prompt.Document) []prompt.Suggest {
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
}

func NewExecutor(repl REPL) prompt.Executor {
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
