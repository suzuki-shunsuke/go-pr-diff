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
	headURL := pr.GetHead().GetRepo().GetCloneURL()

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

	for _, args := range [][]string{
		{"fetch", "--depth", "1", baseURL, baseSHA},
		{"fetch", "--depth", "1", headURL, headSHA},
		{"fetch", "--depth", "1", baseURL, mergeBase},
	} {
		if err := runGit(ctx, dir, args...); err != nil {
			return "", err
		}
	}

	var stdout bytes.Buffer
	cmd := command(ctx, "git", "-C", dir, "diff", baseSHA+"..."+headSHA)
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
