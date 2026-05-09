package prdiff

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v85/github"
)

// Client is a client to get the diff of a pull request.
type Client struct {
	pr   prsService
	repo reposService
}

// NewClient creates a new Client.
// baseURL is a base URL for GitHub Enterprise.
// If baseURL is empty, github.com is used.
// hc is used to call GitHub API using google/go-github.
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
