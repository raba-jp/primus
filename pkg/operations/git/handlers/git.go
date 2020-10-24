package handlers

import (
	"context"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/writer"
	"golang.org/x/xerrors"

	git "github.com/go-git/go-git/v5"
)

type CloneParams struct {
	Path   string
	Branch string
	URL    string
	Cwd    string
}

type CloneHandler interface {
	Clone(ctx context.Context, dryrun bool, p *CloneParams) (err error)
}

type CloneHandlerFunc func(ctx context.Context, dryrun bool, p *CloneParams) error

func (f CloneHandlerFunc) Clone(ctx context.Context, dryrun bool, p *CloneParams) error {
	return f(ctx, dryrun, p)
}

type StorerInitializer func(string) (storage.Storer, billy.Filesystem)

func SetMemoryStore() StorerInitializer {
	return func(_ string) (storage.Storer, billy.Filesystem) { return memory.NewStorage(), nil }
}

func SetFileSystemStore() StorerInitializer {
	return func(path string) (storage.Storer, billy.Filesystem) {
		meta := filepath.Join(path, ".git")
		return filesystem.NewStorage(osfs.New(meta), cache.NewObjectLRUDefault()), osfs.New(path)
	}
}

func NewClone(init StorerInitializer) CloneHandler {
	return CloneHandlerFunc(func(ctx context.Context, dryrun bool, p *CloneParams) error {
		if dryrun {
			if p.Branch != "" {
				ui.Printf("git clone -b %s %s %s\n", p.Branch, p.URL, p.Path)
			} else {
				ui.Printf("git clone %s %s\n", p.URL, p.Path)
			}
			return nil
		}

		base := p.Path
		if !filepath.IsAbs(p.Path) {
			base = filepath.Join(p.Cwd, p.Path)
		}

		storage, fs := init(base)
		if _, err := git.Clone(storage, fs, &git.CloneOptions{
			URL:           p.URL,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + p.Branch),
			Progress:      &writer.NopWriter{},
			SingleBranch:  true,
			Depth:         1,
		}); err != nil {
			return xerrors.Errorf("Failed to git clone: %w", err)
		}

		return nil
	})
}
