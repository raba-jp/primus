//go:generate mockgen -destination mock/handler.go . SetVariableHandler,SetPathHandler

package handlers

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func init() {
	pp.ColoringEnabled = false
}

type VariableScope int

const (
	UniversalScope VariableScope = iota + 1
	GlobalScope
	LocalScope
)

type SetVariableParams struct {
	Name   string
	Value  string
	Scope  VariableScope
	Export bool
}

func (p *SetVariableParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type SetVariableHandler interface {
	SetVariable(ctx context.Context, dryrun bool, p *SetVariableParams) (err error)
}

type SetVariableHandlerFunc func(ctx context.Context, dryrun bool, p *SetVariableParams) error

func (f SetVariableHandlerFunc) SetVariable(ctx context.Context, dryrun bool, p *SetVariableParams) error {
	return f(ctx, dryrun, p)
}

type SetPathParams struct {
	Values []string
}

func (p *SetPathParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type SetPathHandler interface {
	SetPath(ctx context.Context, dryrun bool, p *SetPathParams) (err error)
}

type SetPathHandlerFunc func(ctx context.Context, dryrun bool, p *SetPathParams) error

func (f SetPathHandlerFunc) SetPath(ctx context.Context, dryrun bool, p *SetPathParams) error {
	return f(ctx, dryrun, p)
}

func NewSetVariable(execIF exec.Interface) SetVariableHandler {
	return SetVariableHandlerFunc(func(ctx context.Context, dryrun bool, p *SetVariableParams) error {
		var scope string
		switch p.Scope {
		case UniversalScope:
			scope = "--universal"
		case GlobalScope:
			scope = "--global"
		case LocalScope:
			scope = "--local"
		}

		export := ""
		if p.Export {
			export = " --export"
		}

		arg := fmt.Sprintf("'set %s%s %s %s'", scope, export, p.Name, p.Value)

		if dryrun {
			ui.Printf("fish --command %s\n", arg)
			return nil
		}

		cmd := execIF.CommandContext(ctx, "fish", "--command", arg)
		buf := new(bytes.Buffer)
		errbuf := new(bytes.Buffer)
		cmd.SetStdout(buf)
		cmd.SetStderr(errbuf)
		if err := cmd.Run(); err != nil {
			return xerrors.Errorf("failed to set variable: fish --command %s: %w", arg, err)
		}
		zap.L().Info(
			"set fish variable",
			zap.String("name", p.Name),
			zap.String("value", p.Value),
			zap.String("scope", scope),
			zap.Bool("export", p.Export),
			zap.String("stdout", buf.String()),
			zap.String("stderr", errbuf.String()),
		)
		return nil
	})
}

func NewSetPath(execIF exec.Interface) SetPathHandler {
	return SetPathHandlerFunc(func(ctx context.Context, dryrun bool, p *SetPathParams) error {
		path := fmt.Sprintf("'set --universal fish_user_paths %s'", strings.Join(p.Values, " "))

		if dryrun {
			ui.Printf("fish --command %s\n", path)
			return nil
		}

		cmd := execIF.CommandContext(ctx, "fish", "--command", path)
		buf := new(bytes.Buffer)
		errbuf := new(bytes.Buffer)
		cmd.SetStdout(buf)
		cmd.SetStderr(errbuf)
		if err := cmd.Run(); err != nil {
			return xerrors.Errorf("failed to set path: fish --command 'set --universal fish_user_path %s': %w", path, err)
		}
		zap.L().Info("set fish path", zap.Strings("values", p.Values))
		return nil
	})
}
