package network

import (
	"context"
	"path/filepath"

	lib "go.starlark.net/starlark"

	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type GitCloneParams struct {
	Path   string
	Branch string
	URL    string
	Cwd    string
}

type GitCloneRunner func(ctx context.Context, p *GitCloneParams) error

func NewGitCloneFunction(runner GitCloneRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		params := &GitCloneParams{}
		if err := lib.UnpackArgs(
			b.Name(), args, kwargs,
			"url", &params.URL,
			"path", &params.Path,
			"branch?", &params.Branch,
		); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse argumetns: %w", err)
		}

		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		return lib.None, nil
	}
}

func GitClone(execute command.ExecuteRunner) GitCloneRunner {
	return func(ctx context.Context, p *GitCloneParams) error {
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "git_clone")
		base := p.Path
		if !filepath.IsAbs(p.Path) {
			base = filepath.Join(p.Cwd, p.Path)
		}

		args := []string{"clone", p.URL, base}
		if p.Branch != "" {
			args = []string{"clone", "-b", p.Branch, p.URL, base}
		}

		if err := execute(ctx, &command.Params{Cmd: "git", Args: args}); err != nil {
			return xerrors.Errorf("Failed to git clone: %w", err)
		}

		logger.Info("Finish git clone",
			zap.String("url", p.URL),
			zap.String("path", base),
			zap.String("branch", p.Branch),
		)
		return nil
	}
}
