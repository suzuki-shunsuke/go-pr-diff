package prdiff

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v85/github"
)

type fakePRs struct {
	rawDiff string
	rawErr  error
	pr      *github.PullRequest
	getErr  error
}

func (f *fakePRs) GetRaw(_ context.Context, _, _ string, _ int, _ github.RawOptions) (string, *github.Response, error) {
	return f.rawDiff, nil, f.rawErr
}

func (f *fakePRs) Get(_ context.Context, _, _ string, _ int) (*github.PullRequest, *github.Response, error) {
	return f.pr, nil, f.getErr
}

type fakeRepos struct {
	cmp    *github.CommitsComparison
	cmpErr error
}

func (f *fakeRepos) CompareCommits(_ context.Context, _, _, _, _ string, _ *github.ListOptions) (*github.CommitsComparison, *github.Response, error) {
	return f.cmp, nil, f.cmpErr
}

func TestClient_GetDiff(t *testing.T) {
	t.Parallel()
	wantDiff := "diff --git a/foo b/foo\n"
	apiErr := errors.New("api boom")
	fallbackErr := errors.New("get pr boom")

	tests := []struct {
		name       string
		pr         *fakePRs
		repo       *fakeRepos
		wantDiff   string
		wantErrs   []error
		wantErrSub string
	}{
		{
			name:     "api success returns diff",
			pr:       &fakePRs{rawDiff: wantDiff},
			repo:     &fakeRepos{},
			wantDiff: wantDiff,
		},
		{
			name:     "api fails and fallback fails returns joined error",
			pr:       &fakePRs{rawErr: apiErr, getErr: fallbackErr},
			repo:     &fakeRepos{},
			wantErrs: []error{apiErr, fallbackErr},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := &Client{pr: tc.pr, repo: tc.repo}
			got, err := c.GetDiff(t.Context(), "o", "r", 1)
			if len(tc.wantErrs) == 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if diff := cmp.Diff(tc.wantDiff, got); diff != "" {
					t.Errorf("diff mismatch (-want +got):\n%s", diff)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			for _, we := range tc.wantErrs {
				if !errors.Is(err, we) {
					t.Errorf("error chain missing %v: %v", we, err)
				}
			}
			if tc.wantErrSub != "" && !strings.Contains(err.Error(), tc.wantErrSub) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.wantErrSub)
			}
		})
	}
}

func TestClient_gitFallback_mkghtagPR(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skip network test in -short mode")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not available: %v", err)
	}

	const (
		baseURL   = "https://github.com/suzuki-shunsuke/mkghtag.git"
		headURL   = "https://github.com/ponkio-o/mkghtag.git"
		baseSHA   = "bad96f14de65e0020251d2011221bbce403b2561"
		headSHA   = "6f1db720d939a40b54c472a67cc722edf80d64af"
		mergeBase = baseSHA
	)

	pr := &github.PullRequest{
		Base: &github.PullRequestBranch{
			SHA:  github.Ptr(baseSHA),
			Repo: &github.Repository{CloneURL: github.Ptr(baseURL)},
		},
		Head: &github.PullRequestBranch{
			SHA:  github.Ptr(headSHA),
			Repo: &github.Repository{CloneURL: github.Ptr(headURL)},
		},
	}
	cmpResp := &github.CommitsComparison{
		MergeBaseCommit: &github.RepositoryCommit{SHA: github.Ptr(mergeBase)},
	}

	c := &Client{
		pr:   &fakePRs{pr: pr},
		repo: &fakeRepos{cmp: cmpResp},
	}
	// https://github.com/suzuki-shunsuke/mkghtag/pull/607
	got, err := c.gitFallback(t.Context(), "suzuki-shunsuke", "mkghtag", 607)
	if err != nil {
		t.Fatalf("gitFallback: %v", err)
	}
	want, err := os.ReadFile("testdata/diff_607_git.txt")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if d := cmp.Diff(string(want), got); d != "" {
		t.Errorf("diff mismatch (-want +got):\n%s", d)
	}
}
