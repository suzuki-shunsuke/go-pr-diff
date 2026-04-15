package prdiff

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v84/github"
)

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
