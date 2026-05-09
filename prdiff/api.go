package prdiff

import (
	"context"

	"github.com/google/go-github/v85/github"
)

type reposService interface {
	CompareCommits(ctx context.Context, owner, repo, base, head string, opts *github.ListOptions) (*github.CommitsComparison, *github.Response, error)
}

type prsService interface {
	GetRaw(ctx context.Context, owner, repo string, number int, opts github.RawOptions) (string, *github.Response, error)
	Get(ctx context.Context, owner, repo string, number int) (*github.PullRequest, *github.Response, error)
}
