package repl

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func NewPrompt() {
	buf := new(bytes.Buffer)
	p := prompt.New(Executor(buf), NewCompeter)
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

func Executor(buf *bytes.Buffer) func(s string) {
	return func(s string) {
		eof := false
		thread := &starlark.Thread{}
		s = strings.TrimSpace(s)
		if s == "exit" || s == "quit" {
			os.Exit(0)
			return
		}

		if s == "" {
			return
		}
		buf.WriteString(s + "\n")

		reader := bufio.NewReader(buf)
		f, err := syntax.ParseCompoundStmt("<stdin>", func() ([]byte, error) {
			line, _, err := reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					eof = true
				}
				return nil, err
			}
			return []byte(string(line) + "\n"), nil
		})
		if err != nil {
			if eof {
				return
			}
			ui.Errorf(err.Error())
			return
		}

		if expr := soleExpr(f); expr != nil {
			v, err := starlark.EvalExpr(thread, expr, make(starlark.StringDict))
			if err != nil {
				ui.Errorf(err.Error())
				return
			}
			if v != starlark.None {
				ui.Printf("%v\n", v)
			}
			return
		}
		if err := starlark.ExecREPLChunk(f, thread, nil); err != nil {
			ui.Errorf(err.Error())
			return
		}
		//
		// _, err := starlark.ExecFile(&starlark.Thread{}, "dummy.star", []byte(s), functions.NewPredeclaredFunction(backend.New(exec.New(), afero.NewOsFs()), exec.New(), afero.NewOsFs()))

		// if err != nil {
		// 	panic(err)
		// }
	}
}

func soleExpr(f *syntax.File) syntax.Expr {
	if len(f.Stmts) == 1 {
		if stmt, ok := f.Stmts[0].(*syntax.ExprStmt); ok {
			return stmt.X
		}
	}
	return nil
}
