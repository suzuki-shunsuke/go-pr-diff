package prdiff

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v85/github"
)

// GetDiff returns the diff of a pull request.
// GetDiff first tries to get the diff via the GitHub API, and falls back to git command if it fails.
// If the current directory is a git repository, the pull request base head SHA and merge-base SHA are fetched via git command in the current repository,
// meaning the current repository may be polluted.
func (c *Client) GetDiff(ctx context.Context, owner, repo string, number int) (string, error) {
	diff, _, err := c.pr.GetRaw(ctx, owner, repo, number, github.RawOptions{
		Type: github.Diff,
	})
	if err == nil {
		return diff, nil
	}
	apiErr := fmt.Errorf("get diff via GitHub API: %w", err)
	diff, fbErr := c.gitFallback(ctx, owner, repo, number)
	if fbErr != nil {
		return "", errors.Join(apiErr, fmt.Errorf("get diff via git fallback: %w", fbErr))
	}
	return diff, nil
}
