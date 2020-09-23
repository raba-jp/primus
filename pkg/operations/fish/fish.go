package fish

import (
	"context"

	"github.com/k0kubun/pp"
)

func init() {
	pp.ColoringEnabled = false
}

type VariableScope int

const (
	FishVariableUniversalScope VariableScope = iota + 1
	FishVariableGlobalScope
	FishVariableLocalScope
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
	FishSetVariable(ctx context.Context, dryrun bool, p *SetVariableParams) error
}

type SetVariableHandlerFunc func(ctx context.Context, dryrun bool, p *SetVariableParams) error

func (f SetVariableHandlerFunc) FishSetVariable(ctx context.Context, dryrun bool, p *SetVariableParams) error {
	return f(ctx, dryrun, p)
}

type SetPathParams struct {
	Values []string
}

func (p *SetPathParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type SetPathHandler interface {
	FishSetPath(ctx context.Context, dryrun bool, p *SetPathParams) error
}

type SetPathHandlerFunc func(ctx context.Context, dryrun bool, p *SetPathParams) error

func (f SetPathHandlerFunc) FishSetPath(ctx context.Context, dryrun bool, p *SetPathParams) error {
	return f(ctx, dryrun, p)
}
