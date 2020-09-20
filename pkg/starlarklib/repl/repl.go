package repl

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

type REPL interface {
	Eval(input string) error
}

type repl struct {
	thread      *starlark.Thread
	predeclared functions.Predeclared
	state       *State
}

func NewREPL(state *State, thread *starlark.Thread, predeclared functions.Predeclared) REPL {
	return &repl{
		state:       state,
		thread:      thread,
		predeclared: predeclared,
	}
}

func (r *repl) Eval(input string) error {
	r.state.AppendInput(input)

	f, err := syntax.ParseCompoundStmt("<stdin>", r.state.Readline)
	if err != nil {
		if r.state.Continuation {
			return nil
		}
		return err
	}
	defer r.state.Reset()

	if expr := r.soleExpr(f); expr != nil {
		v, err := starlark.EvalExpr(r.thread, expr, r.predeclared)
		if err != nil {
			return err
		}
		if v != starlark.None {
			ui.Printf("%v\n", v)
		}
		return nil
	}
	if err := starlark.ExecREPLChunk(f, r.thread, r.predeclared); err != nil {
		return err
	}

	return nil
}

func (r *repl) soleExpr(f *syntax.File) syntax.Expr {
	if len(f.Stmts) == 1 {
		if stmt, ok := f.Stmts[0].(*syntax.ExprStmt); ok {
			return stmt.X
		}
	}
	return nil
}

func newThread() *starlark.Thread {
	return &starlark.Thread{
		Name: "REPL",
	}
}
