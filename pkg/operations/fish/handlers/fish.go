//go:generate mockery -outpkg=mocks -case=snake -name=SetVariableHandler
//go:generate mockery -outpkg=mocks -case=snake -name=SetPathHandler

package handlers

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

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

type SetVariableHandler interface {
	Run(ctx context.Context, p *SetVariableParams) (err error)
}

type SetVariableHandlerFunc func(ctx context.Context, p *SetVariableParams) error

func (f SetVariableHandlerFunc) Run(ctx context.Context, p *SetVariableParams) error {
	return f(ctx, p)
}

type SetPathParams struct {
	Values []string
}

type SetPathHandler interface {
	Run(ctx context.Context, p *SetPathParams) (err error)
}

type SetPathHandlerFunc func(ctx context.Context, p *SetPathParams) error

func (f SetPathHandlerFunc) Run(ctx context.Context, p *SetPathParams) error {
	return f(ctx, p)
}

func NewSetVariable(exc exec.Interface) SetVariableHandler {
	return SetVariableHandlerFunc(func(ctx context.Context, p *SetVariableParams) error {
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "fish_set_variable")
		dryrun := ctxlib.DryRun(ctx)

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

		cmd := exc.CommandContext(ctx, "fish", "--command", arg)
		buf := new(bytes.Buffer)
		errbuf := new(bytes.Buffer)
		cmd.SetStdout(buf)
		cmd.SetStderr(errbuf)
		if err := cmd.Run(); err != nil {
			return xerrors.Errorf("failed to set variable: fish --command %s: %w", arg, err)
		}
		logger.Info(
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

func NewSetPath(exc exec.Interface) SetPathHandler {
	return SetPathHandlerFunc(func(ctx context.Context, p *SetPathParams) error {
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "fish_set_path")
		dryrun := ctxlib.DryRun(ctx)
		path := fmt.Sprintf("'set --universal fish_user_paths %s'", strings.Join(p.Values, " "))

		if dryrun {
			ui.Printf("fish --command %s\n", path)
			return nil
		}

		cmd := exc.CommandContext(ctx, "fish", "--command", path)
		buf := new(bytes.Buffer)
		errbuf := new(bytes.Buffer)
		cmd.SetStdout(buf)
		cmd.SetStderr(errbuf)
		if err := cmd.Run(); err != nil {
			return xerrors.Errorf("failed to set path: fish --command 'set --universal fish_user_path %s': %w", path, err)
		}
		logger.Info("set fish path", zap.Strings("values", p.Values))
		return nil
	})
}
