package prdiff

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v84/github"
)

type Client struct {
	pr   prsService
	repo reposService
}

func NewClient(hc *http.Client, baseURL string) (*Client, error) {
	gh := github.NewClient(hc)
	if baseURL != "" {
		g, err := gh.WithEnterpriseURLs(baseURL, "")
		if err != nil {
			return nil, fmt.Errorf("set enterprise URL: %w", err)
		}
		gh = g
	}
	return &Client{
		pr:   gh.PullRequests,
		repo: gh.Repositories,
	}, nil
}
