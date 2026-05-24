package prdiff

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v88/github"
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
	var opts []github.ClientOptionsFunc
	if hc != nil {
		opts = append(opts, github.WithHTTPClient(hc))
	}
	if baseURL != "" {
		opts = append(opts, github.WithEnterpriseURLs(baseURL, baseURL))
	}
	gh, err := github.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("create a GitHub client: %w", err)
	}
	return &Client{
		pr:   gh.PullRequests,
		repo: gh.Repositories,
	}, nil
}
