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
		s := []prompt.Suggest{}

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
