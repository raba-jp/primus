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
		{Text: "command", Description: "Execute command"},
		{Text: "file_copy", Description: "File copy"},
		{Text: "file_move", Description: "File move"},
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
