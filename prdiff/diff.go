package prdiff

import (
	"context"
	"net/http"

	"github.com/google/go-github/v84/github"
)

func NewClient(hc *http.Client, baseURL string) (*Client, error) {
	gh := github.NewClient(hc)
	if baseURL != "" {
		g, err := gh.WithEnterpriseURLs(baseURL, "")
		if err != nil {
			return nil, err
		}
		gh = g
	}
	return &Client{
		pr:   gh.PullRequests,
		repo: gh.Repositories,
	}, nil
}

type reposService interface {
	CompareCommits(ctx context.Context, owner, repo, base, head string, opts *github.ListOptions) (*github.CommitsComparison, *github.Response, error)
}

type prsService interface {
	GetRaw(ctx context.Context, owner, repo string, number int, opts github.RawOptions) (string, *github.Response, error)
	Get(ctx context.Context, owner, repo string, number int) (*github.PullRequest, *github.Response, error)
}

type Client struct {
	pr   prsService
	repo reposService
}

func GetDiff(ctx context.Context, hc *http.Client, owner, repo string, number int) (string, error) {

	return "", nil
}
