package prdiff

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v84/github"
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
