package prdiff

import (
	"bytes"
	"context"
	"fmt"
	"os"
)

func (c *Client) gitFallback(ctx context.Context, owner, repo string, number int) (string, error) {
	pr, _, err := c.pr.Get(ctx, owner, repo, number)
	if err != nil {
		return "", fmt.Errorf("get pull request: %w", err)
	}
	baseSHA := pr.GetBase().GetSHA()
	headSHA := pr.GetHead().GetSHA()
	baseURL := pr.GetBase().GetRepo().GetCloneURL()

	cmp, _, err := c.repo.CompareCommits(ctx, owner, repo, baseSHA, headSHA, nil)
	if err != nil {
		return "", fmt.Errorf("compare commits: %w", err)
	}
	mergeBase := cmp.GetMergeBaseCommit().GetSHA()

	dir, cleanup, err := ensureGitRepo(ctx)
	if err != nil {
		return "", err
	}
	defer cleanup()

	// Always fetch through the base repo URL; GitHub's fork network lets the base
	// repo serve any SHA in the network, so this still works after the head fork
	// is deleted.
	for _, args := range [][]string{
		{"fetch", "--depth", "1", baseURL, mergeBase},
		{"fetch", "--depth", "1", baseURL, headSHA},
	} {
		if err := runGit(ctx, dir, args...); err != nil {
			return "", err
		}
	}

	var stdout bytes.Buffer
	// Use the explicit merge-base instead of `base...head` (3-dot) because the
	// shallow fetches don't have enough history for git to compute merge-base
	// on its own.
	cmd := command(ctx, "git", "-C", dir, "diff", mergeBase, headSHA)
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff: %w: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

func ensureGitRepo(ctx context.Context) (string, func(), error) {
	if isGitRepo(ctx) {
		return ".", func() {}, nil
	}
	dir, err := os.MkdirTemp("", "go-pr-diff-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	if err := runGit(ctx, dir, "init"); err != nil {
		cleanup()
		return "", nil, err
	}
	return dir, cleanup, nil
}

func isGitRepo(ctx context.Context) bool {
	return command(ctx, "git", "rev-parse", "--git-dir").Run() == nil
}

func runGit(ctx context.Context, dir string, args ...string) error {
	full := append([]string{"-C", dir}, args...)
	cmd := command(ctx, "git", full...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %v: %w: %s", args, err, stderr.String())
	}
	return nil
}
